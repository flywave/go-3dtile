package expr

import (
	"fmt"
	"math"
	"strings"
)

type NodeType int

const (
	NODE_LITERAL_NUMBER NodeType = iota
	NODE_LITERAL_STRING
	NODE_LITERAL_BOOLEAN
	NODE_LITERAL_NULL
	NODE_LITERAL_UNDEFINED
	NODE_LITERAL_COLOR
	NODE_LITERAL_VECTOR
	NODE_LITERAL_REGEX
	NODE_VARIABLE
	NODE_BUILTIN_VARIABLE
	NODE_UNARY
	NODE_BINARY
	NODE_TERNARY
	NODE_MEMBER
	NODE_FUNCTION_CALL
	NODE_ARRAY
	NODE_IDENTIFIER
)

type Node struct {
	Type     NodeType
	Token    Token

	NumVal   float64
	StrVal   string
	BoolVal  bool

	Operator string
	Left     *Node
	Right    *Node

	Cond     *Node
	Then     *Node
	Else     *Node

	FuncName string
	Args     []*Node

	Name     string
	Object   *Node

	Elements []*Node

	RegexStr    string
	RegexFlags  string

	VecValues []float64
	Color     [4]float64
}

func (n *Node) String() string {
	if n == nil {
		return "nil"
	}
	switch n.Type {
	case NODE_LITERAL_NUMBER:
		if n.NumVal == math.Trunc(n.NumVal) {
			return fmt.Sprintf("%d", int64(n.NumVal))
		}
		return fmt.Sprintf("%v", n.NumVal)
	case NODE_LITERAL_STRING:
		return fmt.Sprintf("%q", n.StrVal)
	case NODE_LITERAL_BOOLEAN:
		if n.BoolVal {
			return "true"
		}
		return "false"
	case NODE_LITERAL_NULL:
		return "null"
	case NODE_LITERAL_UNDEFINED:
		return "undefined"
	case NODE_LITERAL_COLOR:
		return fmt.Sprintf("color(%v)", n.Color)
	case NODE_LITERAL_VECTOR:
		return fmt.Sprintf("vec%d%v", len(n.VecValues), n.VecValues)
	case NODE_IDENTIFIER:
		return n.Name
	case NODE_VARIABLE:
		return "${" + n.Name + "}"
	case NODE_BUILTIN_VARIABLE:
		return "${" + n.Name + "}"
	case NODE_UNARY:
		return fmt.Sprintf("(%s%s)", n.Operator, n.Left)
	case NODE_BINARY:
		return fmt.Sprintf("(%s %s %s)", n.Left, n.Operator, n.Right)
	case NODE_TERNARY:
		return fmt.Sprintf("(%s ? %s : %s)", n.Cond, n.Then, n.Else)
	case NODE_MEMBER:
		return fmt.Sprintf("(%s.%s)", n.Object, n.Name)
	case NODE_FUNCTION_CALL:
		args := make([]string, len(n.Args))
		for i, a := range n.Args {
			args[i] = a.String()
		}
		return fmt.Sprintf("%s(%s)", n.FuncName, strings.Join(args, ", "))
	case NODE_ARRAY:
		elems := make([]string, len(n.Elements))
		for i, e := range n.Elements {
			elems[i] = e.String()
		}
		return "[" + strings.Join(elems, ", ") + "]"
	case NODE_LITERAL_REGEX:
		return fmt.Sprintf("/%s/%s", n.RegexStr, n.RegexFlags)
	}
	return "?"
}

type ValueType int

const (
	VAL_NUMBER  ValueType = iota
	VAL_STRING
	VAL_BOOLEAN
	VAL_NULL
	VAL_UNDEFINED
	VAL_ARRAY
	VAL_VEC2
	VAL_VEC3
	VAL_VEC4
	VAL_COLOR
	VAL_REGEX
	VAL_OBJECT
)

type Value struct {
	Type    ValueType
	Number  float64
	Str     string
	Boolean bool
	Array   []Value
	Vec     []float64
	Color   [4]float64
	Pattern string
	Flags   string
	Obj     map[string]Value
}

func NumVal(v float64) Value {
	return Value{Type: VAL_NUMBER, Number: v}
}

func StrVal(s string) Value {
	return Value{Type: VAL_STRING, Str: s}
}

func BoolVal(b bool) Value {
	return Value{Type: VAL_BOOLEAN, Boolean: b}
}

func NullVal() Value {
	return Value{Type: VAL_NULL}
}

func UndefinedVal() Value {
	return Value{Type: VAL_UNDEFINED}
}

func ColorVal(r, g, b, a float64) Value {
	return Value{Type: VAL_COLOR, Color: [4]float64{r, g, b, a}}
}

func (v Value) IsTruthy() bool {
	switch v.Type {
	case VAL_NULL, VAL_UNDEFINED:
		return false
	case VAL_BOOLEAN:
		return v.Boolean
	case VAL_NUMBER:
		return v.Number != 0 && !IsNaN(v.Number)
	default:
		return true
	}
}

func (v Value) ToString() string {
	switch v.Type {
	case VAL_NUMBER:
		if v.Number == float64(int64(v.Number)) {
			return fmt.Sprintf("%d", int64(v.Number))
		}
		return fmt.Sprintf("%v", v.Number)
	case VAL_STRING:
		return v.Str
	case VAL_BOOLEAN:
		if v.Boolean {
			return "true"
		}
		return "false"
	case VAL_NULL:
		return "null"
	case VAL_UNDEFINED:
		return "undefined"
	case VAL_COLOR:
		return fmt.Sprintf("(%v,%v,%v,%v)", v.Color[0], v.Color[1], v.Color[2], v.Color[3])
	case VAL_VEC2, VAL_VEC3, VAL_VEC4:
		parts := make([]string, len(v.Vec))
		for i, f := range v.Vec {
			parts[i] = fmt.Sprintf("%v", f)
		}
		return "(" + strings.Join(parts, ",") + ")"
	default:
		return "?"
	}
}

func IsNaN(f float64) bool {
	return f != f
}
