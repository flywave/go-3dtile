package expr

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

type EvalContext struct {
	Features  map[string]Value
	Variables map[string]Value
	Warnings  []string
}

func NewContext() *EvalContext {
	return &EvalContext{
		Features:  make(map[string]Value),
		Variables: make(map[string]Value),
	}
}

func (ctx *EvalContext) Eval(node *Node) (Value, error) {
	if node == nil {
		return NullVal(), nil
	}
	switch node.Type {
	case NODE_LITERAL_NUMBER:
		return NumVal(node.NumVal), nil
	case NODE_LITERAL_STRING:
		return StrVal(node.StrVal), nil
	case NODE_LITERAL_BOOLEAN:
		return BoolVal(node.BoolVal), nil
	case NODE_LITERAL_NULL:
		return NullVal(), nil
	case NODE_LITERAL_UNDEFINED:
		return UndefinedVal(), nil
	case NODE_LITERAL_COLOR:
		return Value{Type: VAL_COLOR, Color: node.Color}, nil
	case NODE_LITERAL_REGEX:
		return Value{Type: VAL_REGEX, Pattern: node.RegexStr, Flags: node.RegexFlags}, nil
	case NODE_IDENTIFIER:
		return ctx.resolveName(node.Name)
	case NODE_VARIABLE:
		return ctx.resolveName(node.Name)
	case NODE_BUILTIN_VARIABLE:
		return ctx.resolveBuiltin(node.Name)
	case NODE_UNARY:
		return ctx.evalUnary(node)
	case NODE_BINARY:
		return ctx.evalBinary(node)
	case NODE_TERNARY:
		return ctx.evalTernary(node)
	case NODE_MEMBER:
		return ctx.evalMember(node)
	case NODE_FUNCTION_CALL:
		return ctx.evalFunction(node)
	case NODE_ARRAY:
		return ctx.evalArray(node)
	case NODE_LITERAL_VECTOR:
		return Value{Type: vecType(len(node.VecValues)), Vec: node.VecValues}, nil
	}
	return UndefinedVal(), nil
}

func vecType(n int) ValueType {
	switch n {
	case 2:
		return VAL_VEC2
	case 3:
		return VAL_VEC3
	case 4:
		return VAL_VEC4
	}
	return VAL_ARRAY
}

func (ctx *EvalContext) resolveName(name string) (Value, error) {
	if v, ok := ctx.Features[name]; ok {
		return v, nil
	}
	if v, ok := ctx.Variables[name]; ok {
		return v, nil
	}
	return UndefinedVal(), nil
}

func (ctx *EvalContext) resolveBuiltin(name string) (Value, error) {
	switch name {
	case "tiles3d_tileset_time":
		return NumVal(0), nil
	}
	return UndefinedVal(), nil
}

func (ctx *EvalContext) evalUnary(n *Node) (Value, error) {
	val, err := ctx.Eval(n.Left)
	if err != nil {
		return UndefinedVal(), err
	}
	switch n.Operator {
	case "-":
		if val.Type == VAL_NUMBER {
			return NumVal(-val.Number), nil
		}
	case "+":
		if val.Type == VAL_NUMBER {
			return NumVal(val.Number), nil
		}
	case "!":
		return BoolVal(!val.IsTruthy()), nil
	case "~":
		if val.Type == VAL_NUMBER {
			return NumVal(float64(^int64(val.Number))), nil
		}
	}
	return UndefinedVal(), nil
}

func (ctx *EvalContext) evalBinary(n *Node) (Value, error) {
	left, err := ctx.Eval(n.Left)
	if err != nil {
		return UndefinedVal(), err
	}
	right, err := ctx.Eval(n.Right)
	if err != nil {
		return UndefinedVal(), err
	}

	switch n.Operator {
	case "+":
		if left.Type == VAL_NUMBER && right.Type == VAL_NUMBER {
			return NumVal(left.Number + right.Number), nil
		}
		return StrVal(left.ToString() + right.ToString()), nil
	case "-":
		if left.Type == VAL_NUMBER && right.Type == VAL_NUMBER {
			return NumVal(left.Number - right.Number), nil
		}
		return UndefinedVal(), nil
	case "*":
		if left.Type == VAL_NUMBER && right.Type == VAL_NUMBER {
			return NumVal(left.Number * right.Number), nil
		}
		return UndefinedVal(), nil
	case "/":
		if left.Type == VAL_NUMBER && right.Type == VAL_NUMBER {
			if right.Number == 0 {
				return UndefinedVal(), nil
			}
			return NumVal(left.Number / right.Number), nil
		}
		return UndefinedVal(), nil
	case "%":
		if left.Type == VAL_NUMBER && right.Type == VAL_NUMBER {
			return NumVal(math.Mod(left.Number, right.Number)), nil
		}
		return UndefinedVal(), nil
	case "===":
		return BoolVal(valuesEqual(left, right)), nil
	case "!==":
		return BoolVal(!valuesEqual(left, right)), nil
	case ">":
		if left.Type == VAL_NUMBER && right.Type == VAL_NUMBER {
			return BoolVal(left.Number > right.Number), nil
		}
		return BoolVal(left.ToString() > right.ToString()), nil
	case ">=":
		if left.Type == VAL_NUMBER && right.Type == VAL_NUMBER {
			return BoolVal(left.Number >= right.Number), nil
		}
		return BoolVal(left.ToString() >= right.ToString()), nil
	case "<":
		if left.Type == VAL_NUMBER && right.Type == VAL_NUMBER {
			return BoolVal(left.Number < right.Number), nil
		}
		return BoolVal(left.ToString() < right.ToString()), nil
	case "<=":
		if left.Type == VAL_NUMBER && right.Type == VAL_NUMBER {
			return BoolVal(left.Number <= right.Number), nil
		}
		return BoolVal(left.ToString() <= right.ToString()), nil
	case "||":
		return BoolVal(left.IsTruthy() || right.IsTruthy()), nil
	case "&&":
		return BoolVal(left.IsTruthy() && right.IsTruthy()), nil
	case "|":
		if left.Type == VAL_NUMBER && right.Type == VAL_NUMBER {
			return NumVal(float64(int64(left.Number) | int64(right.Number))), nil
		}
		return UndefinedVal(), nil
	case "=~":
		return ctx.evalRegexMatch(left, right)
	case "!~":
		v, err := ctx.evalRegexMatch(left, right)
		if err != nil {
			return UndefinedVal(), err
		}
		return BoolVal(!v.Boolean), nil
	}
	return UndefinedVal(), nil
}

func (ctx *EvalContext) evalRegexMatch(text, pattern Value) (Value, error) {
	s := text.ToString()
	p := pattern.Pattern
	if p == "" {
		p = pattern.ToString()
	}
	re, err := regexp.Compile(p)
	if err != nil {
		return BoolVal(false), nil
	}
	return BoolVal(re.MatchString(s)), nil
}

func (ctx *EvalContext) evalTernary(n *Node) (Value, error) {
	cond, err := ctx.Eval(n.Cond)
	if err != nil {
		return UndefinedVal(), err
	}
	if cond.IsTruthy() {
		return ctx.Eval(n.Then)
	}
	return ctx.Eval(n.Else)
}

func (ctx *EvalContext) evalMember(n *Node) (Value, error) {
	obj, err := ctx.Eval(n.Object)
	if err != nil {
		return UndefinedVal(), err
	}
	if obj.Type == VAL_OBJECT && obj.Obj != nil {
		if v, ok := obj.Obj[n.Name]; ok {
			return v, nil
		}
	}
	return UndefinedVal(), nil
}

func (ctx *EvalContext) evalArray(n *Node) (Value, error) {
	arr := make([]Value, len(n.Elements))
	for i, elem := range n.Elements {
		v, err := ctx.Eval(elem)
		if err != nil {
			return UndefinedVal(), err
		}
		arr[i] = v
	}
	return Value{Type: VAL_ARRAY, Array: arr}, nil
}

func (ctx *EvalContext) evalFunction(n *Node) (Value, error) {
	args := make([]Value, len(n.Args))
	for i, a := range n.Args {
		v, err := ctx.Eval(a)
		if err != nil {
			return UndefinedVal(), err
		}
		args[i] = v
	}
	switch n.FuncName {
	case "abs":
		return fnUnary(args, math.Abs)
	case "sqrt":
		return fnUnary(args, math.Sqrt)
	case "cos":
		return fnUnary(args, math.Cos)
	case "sin":
		return fnUnary(args, math.Sin)
	case "tan":
		return fnUnary(args, math.Tan)
	case "acos":
		return fnUnary(args, math.Acos)
	case "asin":
		return fnUnary(args, math.Asin)
	case "atan":
		return fnUnary(args, math.Atan)
	case "radians":
		return fnUnary(args, func(f float64) float64 { return f * math.Pi / 180 })
	case "degrees":
		return fnUnary(args, func(f float64) float64 { return f * 180 / math.Pi })
	case "sign":
		return fnUnary(args, func(f float64) float64 {
			if f > 0 {
				return 1
			} else if f < 0 {
				return -1
			}
			return 0
		})
	case "floor":
		return fnUnary(args, math.Floor)
	case "ceil":
		return fnUnary(args, math.Ceil)
	case "round":
		return fnUnary(args, math.Round)
	case "exp":
		return fnUnary(args, math.Exp)
	case "exp2":
		return fnUnary(args, func(f float64) float64 { return math.Pow(2, f) })
	case "log":
		return fnUnary(args, math.Log)
	case "log2":
		return fnUnary(args, math.Log2)
	case "fract":
		return fnUnary(args, func(f float64) float64 { return f - math.Floor(f) })
	case "length":
		return ctx.fnLength(args)
	case "normalize":
		return ctx.fnNormalize(args)
	case "Boolean":
		return fnUnaryBool(args)
	case "Number":
		return ctx.fnNumber(args)
	case "String":
		if len(args) > 0 {
			return StrVal(args[0].ToString()), nil
		}
		return StrVal(""), nil
	case "isNaN":
		if len(args) > 0 && args[0].Type == VAL_NUMBER {
			return BoolVal(IsNaN(args[0].Number)), nil
		}
		return BoolVal(false), nil
	case "isFinite":
		if len(args) > 0 && args[0].Type == VAL_NUMBER {
			return BoolVal(!math.IsInf(args[0].Number, 0) && !IsNaN(args[0].Number)), nil
		}
		return BoolVal(false), nil
	case "isExactClass":
		return ctx.fnClassCheck(args, true)
	case "isClass":
		return ctx.fnClassCheck(args, false)
	case "atan2":
		return fnBinary(args, math.Atan2)
	case "pow":
		return fnBinary(args, math.Pow)
	case "min":
		return fnBinary(args, math.Min)
	case "max":
		return fnBinary(args, math.Max)
	case "distance":
		return ctx.fnDistance(args)
	case "dot":
		return ctx.fnDot(args)
	case "cross":
		return ctx.fnCross(args)
	case "clamp":
		return ctx.fnClamp(args)
	case "mix":
		return ctx.fnMix(args)
	case "color":
		return ctx.fnColor(args)
	case "rgb":
		return ctx.fnRGB(args)
	case "rgba":
		return ctx.fnRGBA(args)
	case "hsl":
		return ctx.fnHSL(args)
	case "hsla":
		return ctx.fnHSLA(args)
	case "vec2":
		return ctx.fnVec(args, 2)
	case "vec3":
		return ctx.fnVec(args, 3)
	case "vec4":
		return ctx.fnVec(args, 4)
	default:
		return UndefinedVal(), fmt.Errorf("unknown function: %s", n.FuncName)
	}
}

func fnUnary(args []Value, fn func(float64) float64) (Value, error) {
	if len(args) != 1 || args[0].Type != VAL_NUMBER {
		return UndefinedVal(), nil
	}
	return NumVal(fn(args[0].Number)), nil
}

func fnUnaryBool(args []Value) (Value, error) {
	if len(args) != 1 || args[0].Type != VAL_NUMBER {
		return BoolVal(false), nil
	}
	return BoolVal(args[0].Number != 0), nil
}

func fnBinary(args []Value, fn func(float64, float64) float64) (Value, error) {
	if len(args) != 2 || args[0].Type != VAL_NUMBER || args[1].Type != VAL_NUMBER {
		return UndefinedVal(), nil
	}
	return NumVal(fn(args[0].Number, args[1].Number)), nil
}

func (ctx *EvalContext) fnLength(args []Value) (Value, error) {
	if len(args) != 1 {
		return UndefinedVal(), nil
	}
	switch args[0].Type {
	case VAL_VEC2:
		d := args[0].Vec
		return NumVal(math.Sqrt(d[0]*d[0] + d[1]*d[1])), nil
	case VAL_VEC3:
		d := args[0].Vec
		return NumVal(math.Sqrt(d[0]*d[0] + d[1]*d[1] + d[2]*d[2])), nil
	case VAL_VEC4:
		d := args[0].Vec
		return NumVal(math.Sqrt(d[0]*d[0] + d[1]*d[1] + d[2]*d[2] + d[3]*d[3])), nil
	case VAL_NUMBER:
		return NumVal(math.Abs(args[0].Number)), nil
	default:
		return NumVal(float64(len(args[0].ToString()))), nil
	}
}

func (ctx *EvalContext) fnNormalize(args []Value) (Value, error) {
	if len(args) != 1 {
		return UndefinedVal(), nil
	}
	v := args[0]
	l, _ := ctx.fnLength(args)
	if l.Type != VAL_NUMBER || l.Number == 0 {
		return v, nil
	}
	switch v.Type {
	case VAL_VEC2:
		return Value{Type: VAL_VEC2, Vec: []float64{v.Vec[0] / l.Number, v.Vec[1] / l.Number}}, nil
	case VAL_VEC3:
		return Value{Type: VAL_VEC3, Vec: []float64{v.Vec[0] / l.Number, v.Vec[1] / l.Number, v.Vec[2] / l.Number}}, nil
	case VAL_VEC4:
		return Value{Type: VAL_VEC4, Vec: []float64{v.Vec[0] / l.Number, v.Vec[1] / l.Number, v.Vec[2] / l.Number, v.Vec[3] / l.Number}}, nil
	}
	return v, nil
}

func (ctx *EvalContext) fnNumber(args []Value) (Value, error) {
	if len(args) == 0 {
		return NumVal(0), nil
	}
	switch args[0].Type {
	case VAL_NUMBER:
		return args[0], nil
	case VAL_STRING:
		var f float64
		fmt.Sscanf(args[0].Str, "%f", &f)
		return NumVal(f), nil
	case VAL_BOOLEAN:
		if args[0].Boolean {
			return NumVal(1), nil
		}
		return NumVal(0), nil
	}
	return NumVal(0), nil
}

func (ctx *EvalContext) fnClassCheck(args []Value, exact bool) (Value, error) {
	if len(args) != 1 {
		return BoolVal(false), nil
	}
	className := args[0].ToString()
	if v, ok := ctx.Features["_class"]; ok && v.Type == VAL_STRING {
		if exact {
			return BoolVal(v.Str == className), nil
		}
		return BoolVal(strings.Contains(v.Str, className)), nil
	}
	return BoolVal(false), nil
}

func (ctx *EvalContext) fnDistance(args []Value) (Value, error) {
	if len(args) != 2 {
		return UndefinedVal(), nil
	}
	if args[0].Type == VAL_VEC2 && args[1].Type == VAL_VEC2 {
		dx := args[0].Vec[0] - args[1].Vec[0]
		dy := args[0].Vec[1] - args[1].Vec[1]
		return NumVal(math.Sqrt(dx*dx + dy*dy)), nil
	}
	if args[0].Type == VAL_VEC3 && args[1].Type == VAL_VEC3 {
		dx := args[0].Vec[0] - args[1].Vec[0]
		dy := args[0].Vec[1] - args[1].Vec[1]
		dz := args[0].Vec[2] - args[1].Vec[2]
		return NumVal(math.Sqrt(dx*dx + dy*dy + dz*dz)), nil
	}
	n1, _ := ctx.fnNumber([]Value{args[0]})
	n2, _ := ctx.fnNumber([]Value{args[1]})
	return NumVal(math.Abs(n1.Number - n2.Number)), nil
}

func (ctx *EvalContext) fnDot(args []Value) (Value, error) {
	if len(args) != 2 {
		return UndefinedVal(), nil
	}
	if args[0].Type == VAL_VEC3 && args[1].Type == VAL_VEC3 {
		d := args[0].Vec[0]*args[1].Vec[0] + args[0].Vec[1]*args[1].Vec[1] + args[0].Vec[2]*args[1].Vec[2]
		return NumVal(d), nil
	}
	if args[0].Type == VAL_VEC2 && args[1].Type == VAL_VEC2 {
		d := args[0].Vec[0]*args[1].Vec[0] + args[0].Vec[1]*args[1].Vec[1]
		return NumVal(d), nil
	}
	return UndefinedVal(), nil
}

func (ctx *EvalContext) fnCross(args []Value) (Value, error) {
	if len(args) != 2 || args[0].Type != VAL_VEC3 || args[1].Type != VAL_VEC3 {
		return UndefinedVal(), nil
	}
	a, b := args[0].Vec, args[1].Vec
	return Value{Type: VAL_VEC3, Vec: []float64{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}}, nil
}

func (ctx *EvalContext) fnClamp(args []Value) (Value, error) {
	if len(args) != 3 || args[0].Type != VAL_NUMBER || args[1].Type != VAL_NUMBER || args[2].Type != VAL_NUMBER {
		return UndefinedVal(), nil
	}
	return NumVal(math.Max(args[1].Number, math.Min(args[0].Number, args[2].Number))), nil
}

func (ctx *EvalContext) fnMix(args []Value) (Value, error) {
	if len(args) != 3 || args[0].Type != VAL_NUMBER || args[1].Type != VAL_NUMBER || args[2].Type != VAL_NUMBER {
		return UndefinedVal(), nil
	}
	t := args[2].Number
	return NumVal(args[0].Number*(1-t) + args[1].Number*t), nil
}

func (ctx *EvalContext) fnColor(args []Value) (Value, error) {
	if len(args) == 0 {
		return ColorVal(1, 1, 1, 1), nil
	}
	if args[0].Type == VAL_STRING {
		return parseColorString(args[0].Str)
	}
	if args[0].Type == VAL_COLOR {
		return args[0], nil
	}
	if args[0].Type == VAL_ARRAY && len(args[0].Array) >= 3 {
		a := args[0].Array
		r, _ := toFloat64(a[0])
		g, _ := toFloat64(a[1])
		b, _ := toFloat64(a[2])
		alpha := 1.0
		if len(a) >= 4 {
			alpha, _ = toFloat64(a[3])
		}
		return ColorVal(r, g, b, alpha), nil
	}
	return ColorVal(1, 1, 1, 1), nil
}

func (ctx *EvalContext) fnRGB(args []Value) (Value, error) {
	if len(args) >= 3 {
		r, _ := toFloat64(args[0])
		g, _ := toFloat64(args[1])
		b, _ := toFloat64(args[2])
		return ColorVal(r, g, b, 1), nil
	}
	return ColorVal(1, 1, 1, 1), nil
}

func (ctx *EvalContext) fnRGBA(args []Value) (Value, error) {
	if len(args) >= 4 {
		r, _ := toFloat64(args[0])
		g, _ := toFloat64(args[1])
		b, _ := toFloat64(args[2])
		a, _ := toFloat64(args[3])
		return ColorVal(r, g, b, a), nil
	}
	return ColorVal(1, 1, 1, 1), nil
}

func (ctx *EvalContext) fnHSL(args []Value) (Value, error) {
	if len(args) >= 3 {
		h, _ := toFloat64(args[0])
		s, _ := toFloat64(args[1])
		l, _ := toFloat64(args[2])
		c := hslToRGB(h, s, l)
		return ColorVal(c[0], c[1], c[2], 1), nil
	}
	return ColorVal(1, 1, 1, 1), nil
}

func (ctx *EvalContext) fnHSLA(args []Value) (Value, error) {
	if len(args) >= 4 {
		h, _ := toFloat64(args[0])
		s, _ := toFloat64(args[1])
		l, _ := toFloat64(args[2])
		a, _ := toFloat64(args[3])
		c := hslToRGB(h, s, l)
		return ColorVal(c[0], c[1], c[2], a), nil
	}
	return ColorVal(1, 1, 1, 1), nil
}

func (ctx *EvalContext) fnVec(args []Value, dim int) (Value, error) {
	vec := make([]float64, dim)
	for i := 0; i < dim && i < len(args); i++ {
		v, _ := toFloat64(args[i])
		vec[i] = v
	}
	return Value{Type: vecType(dim), Vec: vec}, nil
}

func toFloat64(v Value) (float64, bool) {
	switch v.Type {
	case VAL_NUMBER:
		return v.Number, true
	case VAL_STRING:
		var f float64
		fmt.Sscanf(v.Str, "%f", &f)
		return f, true
	case VAL_BOOLEAN:
		if v.Boolean {
			return 1, true
		}
		return 0, true
	}
	return 0, false
}

func valuesEqual(a, b Value) bool {
	if a.Type != b.Type {
		return false
	}
	switch a.Type {
	case VAL_NUMBER:
		return a.Number == b.Number
	case VAL_STRING:
		return a.Str == b.Str
	case VAL_BOOLEAN:
		return a.Boolean == b.Boolean
	case VAL_NULL, VAL_UNDEFINED:
		return true
	case VAL_COLOR:
		return a.Color == b.Color
	case VAL_ARRAY:
		if len(a.Array) != len(b.Array) {
			return false
		}
		for i := range a.Array {
			if !valuesEqual(a.Array[i], b.Array[i]) {
				return false
			}
		}
		return true
	default:
		return a.ToString() == b.ToString()
	}
}

func parseColorString(s string) (Value, error) {
	switch strings.ToUpper(s) {
	case "WHITE":
		return ColorVal(1, 1, 1, 1), nil
	case "BLACK":
		return ColorVal(0, 0, 0, 1), nil
	case "RED":
		return ColorVal(1, 0, 0, 1), nil
	case "GREEN":
		return ColorVal(0, 1, 0, 1), nil
	case "BLUE":
		return ColorVal(0, 0, 1, 1), nil
	case "YELLOW":
		return ColorVal(1, 1, 0, 1), nil
	case "CYAN":
		return ColorVal(0, 1, 1, 1), nil
	case "MAGENTA":
		return ColorVal(1, 0, 1, 1), nil
	case "TRANSPARENT":
		return ColorVal(0, 0, 0, 0), nil
	}
	if len(s) > 0 && s[0] == '#' {
		return parseHexColor(s)
	}
	return ColorVal(1, 1, 1, 1), nil
}

func parseHexColor(s string) (Value, error) {
	s = s[1:]
	var r, g, b, a uint8
	a = 255
	switch len(s) {
	case 6:
		fmt.Sscanf(s, "%02x%02x%02x", &r, &g, &b)
	case 8:
		fmt.Sscanf(s, "%02x%02x%02x%02x", &r, &g, &b, &a)
	case 3:
		fmt.Sscanf(s, "%1x%1x%1x", &r, &g, &b)
		r, g, b = r*17, g*17, b*17
	default:
		return ColorVal(1, 1, 1, 1), nil
	}
	return ColorVal(float64(r)/255, float64(g)/255, float64(b)/255, float64(a)/255), nil
}

func hslToRGB(h, s, l float64) [4]float64 {
	h = h / 360
	s = s / 100
	l = l / 100
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h*6, 2)-1))
	m := l - c/2

	var r, g, b float64
	switch {
	case h < 1.0/6:
		r, g, b = c, x, 0
	case h < 2.0/6:
		r, g, b = x, c, 0
	case h < 3.0/6:
		r, g, b = 0, c, x
	case h < 4.0/6:
		r, g, b = 0, x, c
	case h < 5.0/6:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	return [4]float64{r + m, g + m, b + m, 1}
}
