package expr

import (
	"math"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input string
		count int
	}{
		{"true", 2},
		{"123", 2},
		{"1.5", 2},
		{"\"hello\"", 2},
		{"${Height}", 4},
		{"1 + 2", 4},
		{"a === b", 4},
		{"/pattern/", 2},
	}
	for _, tt := range tests {
		tokens, err := Lex(tt.input)
		if err != nil {
			t.Errorf("Lex(%q) error: %v", tt.input, err)
			continue
		}
		if len(tokens) != tt.count {
			t.Errorf("Lex(%q) got %d tokens, want %d: %v", tt.input, len(tokens), tt.count, tokens)
		}
	}
}

func TestParserBasic(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"true", "true"},
		{"false", "false"},
		{"null", "null"},
		{"42", "42"},
		{"1.5", "1.5"},
		{"\"hello\"", "\"hello\""},
		{"${Height}", "${Height}"},
		{"!x", "(!x)"},
		{"1 + 2", "(1 + 2)"},
		{"1 + 2 * 3", "(1 + (2 * 3))"},
		{"(1 + 2) * 3", "((1 + 2) * 3)"},
		{"a === b", "(a === b)"},
		{"a !== b", "(a !== b)"},
		{"a > b", "(a > b)"},
		{"a >= b", "(a >= b)"},
		{"a < b", "(a < b)"},
		{"a <= b", "(a <= b)"},
		{"a && b", "(a && b)"},
		{"a || b", "(a || b)"},
		{"a ? b : c", "(a ? b : c)"},
		{"a.b", "(a.b)"},
		{"sin(x)", "sin(x)"},
		{"sin(x + 1)", "sin((x + 1))"},
		{"[1, 2, 3]", "[1, 2, 3]"},
		{"a | b", "(a | b)"},
	}
	for _, tt := range tests {
		node, err := Parse(tt.input)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", tt.input, err)
			continue
		}
		got := node.String()
		if got != tt.want {
			t.Errorf("Parse(%q).String() = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestEvalLiteral(t *testing.T) {
	ctx := NewContext()
	v, err := ctx.Eval(&Node{Type: NODE_LITERAL_NUMBER, NumVal: 42})
	if err != nil || v.Number != 42 || v.Type != VAL_NUMBER {
		t.Errorf("number literal failed: %v %v", v, err)
	}
	v, err = ctx.Eval(&Node{Type: NODE_LITERAL_BOOLEAN, BoolVal: true})
	if err != nil || v.Boolean != true {
		t.Errorf("boolean literal failed: %v", v)
	}
	v, err = ctx.Eval(&Node{Type: NODE_LITERAL_STRING, StrVal: "hello"})
	if err != nil || v.Str != "hello" {
		t.Errorf("string literal failed: %v", v)
	}
}

func TestEvalBinary(t *testing.T) {
	ctx := NewContext()
	cases := []struct {
		expr string
		want float64
	}{
		{"1 + 2", 3},
		{"10 - 3", 7},
		{"3 * 4", 12},
		{"10 / 2", 5},
		{"10 % 3", 1},
	}
	for _, c := range cases {
		node, err := Parse(c.expr)
		if err != nil {
			t.Fatalf("Parse(%q) error: %v", c.expr, err)
		}
		v, err := ctx.Eval(node)
		if err != nil || v.Type != VAL_NUMBER || v.Number != c.want {
			t.Errorf("Eval(%q) = %v (err=%v), want %v", c.expr, v.Number, err, c.want)
		}
	}
}

func TestEvalComparison(t *testing.T) {
	if v, _ := NewContext().EvalNode("1 === 1"); v.Boolean != true {
		t.Error("1 === 1 should be true")
	}
	if v, _ := NewContext().EvalNode("1 !== 2"); v.Boolean != true {
		t.Error("1 !== 2 should be true")
	}
	if v, _ := NewContext().EvalNode("2 > 1"); v.Boolean != true {
		t.Error("2 > 1 should be true")
	}
	if v, _ := NewContext().EvalNode("2 < 1"); v.Boolean != false {
		t.Error("2 < 1 should be false")
	}
}

func TestEvalLogical(t *testing.T) {
	if v, _ := NewContext().EvalNode("true && false"); v.Boolean != false {
		t.Error("true && false should be false")
	}
	if v, _ := NewContext().EvalNode("true || false"); v.Boolean != true {
		t.Error("true || false should be true")
	}
	if v, _ := NewContext().EvalNode("!true"); v.Boolean != false {
		t.Error("!true should be false")
	}
}

func TestEvalTernary(t *testing.T) {
	if v, _ := NewContext().EvalNode("true ? 1 : 2"); v.Number != 1 {
		t.Error("true ? 1 : 2 should be 1")
	}
	if v, _ := NewContext().EvalNode("false ? 1 : 2"); v.Number != 2 {
		t.Error("false ? 1 : 2 should be 2")
	}
}

func TestEvalPrecedence(t *testing.T) {
	if v, _ := NewContext().EvalNode("1 + 2 * 3"); v.Number != 7 {
		t.Errorf("1 + 2 * 3 = %v, want 7", v.Number)
	}
	if v, _ := NewContext().EvalNode("(1 + 2) * 3"); v.Number != 9 {
		t.Errorf("(1 + 2) * 3 = %v, want 9", v.Number)
	}
}

func TestEvalFeature(t *testing.T) {
	ctx := NewContext()
	ctx.Features["Height"] = NumVal(100)
	if v, _ := ctx.EvalNode("${Height}"); v.Number != 100 {
		t.Error("${Height} should be 100")
	}
	if v, _ := ctx.EvalNode("${Height} > 50"); v.Boolean != true {
		t.Error("${Height} > 50 should be true")
	}
}

func TestEvalString(t *testing.T) {
	if v, _ := NewContext().EvalNode("\"hello\" + \" world\""); v.Str != "hello world" {
		t.Errorf("string concat = %q, want %q", v.Str, "hello world")
	}
}

func TestEvalArray(t *testing.T) {
	v, err := NewContext().EvalNode("[1, 2, 3]")
	if err != nil || v.Type != VAL_ARRAY || len(v.Array) != 3 {
		t.Errorf("array failed: %v %v", v, err)
	}
}

func TestEvalMathFunc(t *testing.T) {
	tests := []struct {
		expr string
		want float64
	}{
		{"abs(-5)", 5},
		{"sqrt(9)", 3},
		{"floor(3.7)", 3},
		{"ceil(3.2)", 4},
		{"round(3.5)", 4},
		{"min(3, 7)", 3},
		{"max(3, 7)", 7},
		{"pow(2, 3)", 8},
		{"sign(-5)", -1},
		{"sign(5)", 1},
	}
	for _, tt := range tests {
		v, err := NewContext().EvalNode(tt.expr)
		if err != nil || v.Type != VAL_NUMBER || v.Number != tt.want {
			t.Errorf("Eval(%q) = %v (err=%v), want %v", tt.expr, v.Number, err, tt.want)
		}
	}
}

func TestEvalColor(t *testing.T) {
	v, err := NewContext().EvalNode("color('#FF0000')")
	if err != nil || v.Type != VAL_COLOR {
		t.Errorf("color('#FF0000') failed: %v %v", v, err)
	}
	if v.Color[0] != 1 || v.Color[1] != 0 || v.Color[2] != 0 || v.Color[3] != 1 {
		t.Errorf("color('#FF0000') = %v, want (1,0,0,1)", v.Color)
	}

	v, err = NewContext().EvalNode("rgb(1, 0, 0)")
	if err != nil || v.Type != VAL_COLOR {
		t.Errorf("rgb(1,0,0) failed: %v", v)
	}
}

func TestEvalVector(t *testing.T) {
	v, err := NewContext().EvalNode("vec3(1, 2, 3)")
	if err != nil || v.Type != VAL_VEC3 || len(v.Vec) != 3 {
		t.Errorf("vec3 failed: %v %v", v, err)
	}
	if v.Vec[1] != 2 {
		t.Errorf("vec3[1] = %v, want 2", v.Vec[1])
	}
}

func TestEvalMember(t *testing.T) {
	ctx := NewContext()
	ctx.Features["building"] = Value{Type: VAL_OBJECT, Obj: map[string]Value{
		"height": NumVal(50),
	}}
	v, err := ctx.EvalNode("building.height")
	if err != nil || v.Type != VAL_NUMBER || v.Number != 50 {
		t.Errorf("member access failed: %v %v", v, err)
	}
}

func TestPrecedenceChain(t *testing.T) {
	ctx := NewContext()
	ctx.Features["x"] = NumVal(10)
	ctx.Features["y"] = NumVal(20)
	v, err := ctx.EvalNode("${x} + ${y} * 2")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if v.Number != 50 {
		t.Errorf("x + y * 2 = %v, want 50", v.Number)
	}
}

func TestStringConcat(t *testing.T) {
	v, err := NewContext().EvalNode("\"a\" + \"b\"")
	if err != nil || v.Str != "ab" {
		t.Errorf("string concat = %q, want %q", v.Str, "ab")
	}
}

func TestRegex(t *testing.T) {
	v, err := NewContext().EvalNode("\"hello\" =~ /^h/")
	if err != nil || v.Boolean != true {
		t.Errorf("regex match failed: %v %v", v, err)
	}
	v, err = NewContext().EvalNode("\"world\" =~ /^h/")
	if err != nil || v.Boolean != false {
		t.Errorf("regex non-match failed: %v %v", v, err)
	}
}

func TestComplexExpression(t *testing.T) {
	ctx := NewContext()
	ctx.Features["Height"] = NumVal(100)
	ctx.Features["Area"] = NumVal(500)
	v, err := ctx.EvalNode("(${Height} > 50 && ${Area} > 100) ? color('#FF0000') : color('#FFFFFF')")
	if err != nil {
		t.Fatalf("complex expression error: %v", err)
	}
	if v.Type != VAL_COLOR || v.Color[0] != 1 || v.Color[1] != 0 {
		t.Errorf("complex expression = %v, want red", v)
	}
}

func TestEvalError(t *testing.T) {
	_, err := Parse("1 + ")
	if err == nil {
		t.Error("expected error for incomplete expression")
	}
}

func (ctx *EvalContext) EvalNode(input string) (Value, error) {
	node, err := Parse(input)
	if err != nil {
		return UndefinedVal(), err
	}
	return ctx.Eval(node)
}

func TestEvalTrig(t *testing.T) {
	v, _ := NewContext().EvalNode("sin(0)")
	if math.Abs(v.Number) > 0.0001 {
		t.Errorf("sin(0) = %v, want 0", v.Number)
	}
	v, _ = NewContext().EvalNode("cos(0)")
	if math.Abs(v.Number-1) > 0.0001 {
		t.Errorf("cos(0) = %v, want 1", v.Number)
	}
}

func TestEvalHSL(t *testing.T) {
	v, err := NewContext().EvalNode("hsl(0, 100, 50)")
	if err != nil || v.Type != VAL_COLOR {
		t.Errorf("hsl failed: %v %v", v, err)
	}
	if math.Abs(v.Color[0]-1) > 0.01 || math.Abs(v.Color[1]) > 0.01 || math.Abs(v.Color[2]) > 0.01 {
		t.Errorf("hsl(0,100,50) = %v, want (1,0,0,1)", v.Color)
	}
}

func TestEvalDistance(t *testing.T) {
	v, err := NewContext().EvalNode("distance(vec3(0,0,0), vec3(3,4,0))")
	if err != nil || v.Type != VAL_NUMBER || math.Abs(v.Number-5) > 0.0001 {
		t.Errorf("distance = %v, want 5 (err=%v)", v.Number, err)
	}
}

func TestEvalDot(t *testing.T) {
	v, err := NewContext().EvalNode("dot(vec3(1,0,0), vec3(0,1,0))")
	if err != nil || v.Type != VAL_NUMBER || math.Abs(v.Number) > 0.0001 {
		t.Errorf("dot = %v, want 0", v.Number)
	}
}

func TestEvalCross(t *testing.T) {
	v, err := NewContext().EvalNode("cross(vec3(1,0,0), vec3(0,1,0))")
	if err != nil || v.Type != VAL_VEC3 || math.Abs(v.Vec[2]-1) > 0.0001 {
		t.Errorf("cross = %v, want (0,0,1)", v.Vec)
	}
}

func TestEvalNormalize(t *testing.T) {
	v, err := NewContext().EvalNode("normalize(vec3(3, 4, 0))")
	if err != nil || v.Type != VAL_VEC3 || math.Abs(v.Vec[0]-0.6) > 0.0001 {
		t.Errorf("normalize = %v, want (0.6, 0.8, 0)", v.Vec)
	}
}

func TestEvalClamp(t *testing.T) {
	v, _ := NewContext().EvalNode("clamp(5, 0, 3)")
	if v.Number != 3 {
		t.Errorf("clamp(5,0,3) = %v, want 3", v.Number)
	}
	v, _ = NewContext().EvalNode("clamp(-1, 0, 3)")
	if v.Number != 0 {
		t.Errorf("clamp(-1,0,3) = %v, want 0", v.Number)
	}
}

func TestEvalBitwise(t *testing.T) {
	v, _ := NewContext().EvalNode("1 | 2")
	if v.Number != 3 {
		t.Errorf("1 | 2 = %v, want 3", v.Number)
	}
}

func TestEvalStringFuncs(t *testing.T) {
	v, _ := NewContext().EvalNode("Number('42')")
	if v.Number != 42 {
		t.Errorf("Number('42') = %v, want 42", v.Number)
	}
	v, _ = NewContext().EvalNode("String(42)")
	if v.Str != "42" {
		t.Errorf("String(42) = %q, want 42", v.Str)
	}
}

func TestEvalNegate(t *testing.T) {
	v, _ := NewContext().EvalNode("-5")
	if v.Number != -5 {
		t.Errorf("-5 = %v, want -5", v.Number)
	}
}

func TestEvalRegexNotMatch(t *testing.T) {
	v, err := NewContext().EvalNode("\"hello\" !~ /^x/")
	if err != nil || v.Boolean != true {
		t.Errorf("regex not-match should be true: %v %v", v, err)
	}
}
