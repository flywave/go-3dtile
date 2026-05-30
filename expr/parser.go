package expr

import "fmt"

type Parser struct {
	tokens []Token
	pos    int
}

func Parse(input string) (*Node, error) {
	tokens, err := Lex(input)
	if err != nil {
		return nil, fmt.Errorf("lex error: %w", err)
	}
	p := &Parser{tokens: tokens}
	node := p.parseExpression()
	if p.err() != nil {
		return nil, p.err()
	}
	return node, nil
}

func (p *Parser) peek() Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return Token{Type: TT_EOF}
}

func (p *Parser) advance() Token {
	t := p.peek()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return t
}

func (p *Parser) expect(tt TokenType) (Token, error) {
	if p.peek().Type == tt {
		return p.advance(), nil
	}
	return Token{}, fmt.Errorf("expected %d at position %d, got %s", tt, p.peek().Pos, p.peek().Value)
}

func (p *Parser) match(types ...TokenType) bool {
	t := p.peek().Type
	for _, tt := range types {
		if t == tt {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) err() error {
	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == TT_INVALID {
		return fmt.Errorf("invalid token at %d: %q", p.tokens[p.pos].Pos, p.tokens[p.pos].Value)
	}
	return nil
}

// expression → conditional
func (p *Parser) parseExpression() *Node {
	return p.parseConditional()
}

// conditional → logical_or ("?" expression ":" expression)?
func (p *Parser) parseConditional() *Node {
	node := p.parseLogicalOr()
	if p.match(TT_QUESTION) {
		thenBranch := p.parseExpression()
		p.expect(TT_COLON)
		elseBranch := p.parseExpression()
		node = &Node{Type: NODE_TERNARY, Cond: node, Then: thenBranch, Else: elseBranch}
	}
	return node
}

// logical_or → logical_and ("||" logical_and)*
func (p *Parser) parseLogicalOr() *Node {
	node := p.parseLogicalAnd()
	for p.peek().Type == TT_OR {
		p.advance()
		right := p.parseLogicalAnd()
		node = &Node{Type: NODE_BINARY, Operator: "||", Left: node, Right: right}
	}
	return node
}

// logical_and → equality ("&&" equality)*
func (p *Parser) parseLogicalAnd() *Node {
	node := p.parseBitwiseOr()
	for p.peek().Type == TT_AND {
		p.advance()
		right := p.parseBitwiseOr()
		node = &Node{Type: NODE_BINARY, Operator: "&&", Left: node, Right: right}
	}
	return node
}

// bitwise_or → equality ("|" equality)*
func (p *Parser) parseBitwiseOr() *Node {
	node := p.parseEquality()
	for p.peek().Type == TT_BAR {
		p.advance()
		right := p.parseEquality()
		node = &Node{Type: NODE_BINARY, Operator: "|", Left: node, Right: right}
	}
	return node
}

// equality → relational (("===" | "!==" | "=~" | "!~") relational)*
func (p *Parser) parseEquality() *Node {
	node := p.parseRelational()
	for p.peek().Type == TT_EQ || p.peek().Type == TT_NEQ || p.peek().Type == TT_REGEX_MATCH || p.peek().Type == TT_REGEX_NOT {
		op := p.advance().Value
		right := p.parseRelational()
		node = &Node{Type: NODE_BINARY, Operator: op, Left: node, Right: right}
	}
	return node
}

// relational → additive ((">" | ">=" | "<" | "<=") additive)*
func (p *Parser) parseRelational() *Node {
	node := p.parseAdditive()
	for p.peek().Type == TT_GT || p.peek().Type == TT_GTE || p.peek().Type == TT_LT || p.peek().Type == TT_LTE {
		op := p.advance().Value
		right := p.parseAdditive()
		node = &Node{Type: NODE_BINARY, Operator: op, Left: node, Right: right}
	}
	return node
}

// additive → multiplicative (("+" | "-") multiplicative)*
func (p *Parser) parseAdditive() *Node {
	node := p.parseMultiplicative()
	for p.peek().Type == TT_PLUS || p.peek().Type == TT_MINUS {
		op := p.advance().Value
		right := p.parseMultiplicative()
		node = &Node{Type: NODE_BINARY, Operator: op, Left: node, Right: right}
	}
	return node
}

// multiplicative → unary (("*" | "/" | "%") unary)*
func (p *Parser) parseMultiplicative() *Node {
	node := p.parseUnary()
	for p.peek().Type == TT_STAR || p.peek().Type == TT_SLASH || p.peek().Type == TT_PERCENT {
		op := p.advance().Value
		right := p.parseUnary()
		node = &Node{Type: NODE_BINARY, Operator: op, Left: node, Right: right}
	}
	return node
}

// unary → ("!" | "-" | "+" | "~") unary | postfix
func (p *Parser) parseUnary() *Node {
	if p.peek().Type == TT_BANG || p.peek().Type == TT_MINUS || p.peek().Type == TT_PLUS || p.peek().Type == TT_TILDE {
		op := p.advance().Value
		right := p.parseUnary()
		return &Node{Type: NODE_UNARY, Operator: op, Left: right}
	}
	return p.parsePostfix()
}

// postfix → primary ("." IDENTIFIER)*
func (p *Parser) parsePostfix() *Node {
	node := p.parsePrimary()
	for p.peek().Type == TT_DOT {
		p.advance() // consume .
		tok, err := p.expect(TT_IDENTIFIER)
		if err != nil {
			return node
		}
		node = &Node{Type: NODE_MEMBER, Object: node, Name: tok.Value}
	}
	return node
}

// primary → NUMBER | STRING | BOOLEAN | "null" | "undefined"
//
//	| "[" list "]"
//	| IDENTIFIER "(" args ")"
//	| "${" expression "}"
//	| "/" pattern "/" flags
//	| "(" expression ")"
//	| IDENTIFIER   (variables, builtins, or color/vector constructors)
func (p *Parser) parsePrimary() *Node {
	tok := p.peek()

	switch tok.Type {
	case TT_NUMBER:
		p.advance()
		return p.makeNumber(tok.Value)

	case TT_STRING:
		p.advance()
		return &Node{Type: NODE_LITERAL_STRING, StrVal: tok.Value}

	case TT_TRUE:
		p.advance()
		return &Node{Type: NODE_LITERAL_BOOLEAN, BoolVal: true}

	case TT_FALSE:
		p.advance()
		return &Node{Type: NODE_LITERAL_BOOLEAN, BoolVal: false}

	case TT_NULL:
		p.advance()
		return &Node{Type: NODE_LITERAL_NULL}

	case TT_UNDEFINED:
		p.advance()
		return &Node{Type: NODE_LITERAL_UNDEFINED}

	case TT_LBRACKET:
		p.advance()
		var elems []*Node
		if p.peek().Type != TT_RBRACKET {
			elems = append(elems, p.parseExpression())
			for p.match(TT_COMMA) {
				elems = append(elems, p.parseExpression())
			}
		}
		p.expect(TT_RBRACKET)
		return &Node{Type: NODE_ARRAY, Elements: elems}

	case TT_LPAREN:
		p.advance()
		node := p.parseExpression()
		p.expect(TT_RPAREN)
		return node

	case TT_VARIABLE_OPEN:
		p.advance()
		node := p.parseExpression()
		if _, err := p.expect(TT_VARIABLE_CLOSE); err != nil {
			p.error(err.Error())
			return node
		}
		node.Type = NODE_VARIABLE
		return node

	case TT_REGEX:
		p.advance()
		return p.parseRegexLiteral(tok.Value)

	case TT_IDENTIFIER:
		p.advance()
		name := tok.Value
		if p.peek().Type == TT_LPAREN {
			return p.parseFunctionCall(name)
		}
		// Check for color/vec literal constructors without parens (legacy)
		if c, ok := parseColorConstructor(name); ok {
			return &Node{Type: NODE_LITERAL_COLOR, Color: c}
		}
		return &Node{Type: NODE_IDENTIFIER, Name: name}

	default:
		p.error(fmt.Sprintf("unexpected token %s", tok.Value))
		return &Node{Type: NODE_LITERAL_NULL}
	}
}

func (p *Parser) parseFunctionCall(name string) *Node {
	p.advance() // consume (
	var args []*Node
	if p.peek().Type != TT_RPAREN {
		args = append(args, p.parseExpression())
		for p.match(TT_COMMA) {
			args = append(args, p.parseExpression())
		}
	}
	p.expect(TT_RPAREN)
	return &Node{Type: NODE_FUNCTION_CALL, FuncName: name, Args: args}
}

func (p *Parser) parseRegexLiteral(raw string) *Node {
	// raw looks like /pattern/flags
	s := raw
	flags := ""
	if len(s) > 0 && s[len(s)-1] == '/' {
		s = s[1 : len(s)-1]
	} else if idx := stringsLastIndex(s, "/"); idx > 0 {
		flags = s[idx+1:]
		s = s[1:idx]
	}
	return &Node{Type: NODE_LITERAL_REGEX, RegexStr: s, RegexFlags: flags}
}

func (p *Parser) makeNumber(raw string) *Node {
	var f float64
	fmt.Sscanf(raw, "%f", &f)
	return &Node{Type: NODE_LITERAL_NUMBER, NumVal: f}
}

func (p *Parser) error(msg string) {
	// Store error via the token stream
	if p.pos < len(p.tokens) {
		p.tokens[p.pos].Type = TT_INVALID
		p.tokens[p.pos].Value = msg
	}
}

func parseColorConstructor(name string) ([4]float64, bool) {
	switch name {
	case "WHITE":
		return [4]float64{1, 1, 1, 1}, true
	case "BLACK":
		return [4]float64{0, 0, 0, 1}, true
	case "RED":
		return [4]float64{1, 0, 0, 1}, true
	case "GREEN":
		return [4]float64{0, 1, 0, 1}, true
	case "BLUE":
		return [4]float64{0, 0, 1, 1}, true
	case "YELLOW":
		return [4]float64{1, 1, 0, 1}, true
	case "CYAN":
		return [4]float64{0, 1, 1, 1}, true
	case "MAGENTA":
		return [4]float64{1, 0, 1, 1}, true
	case "TRANSPARENT":
		return [4]float64{0, 0, 0, 0}, true
	}
	return [4]float64{}, false
}

func stringsLastIndex(s, substr string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == substr[0] && i+len(substr) <= len(s) && s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
