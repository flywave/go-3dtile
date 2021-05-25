package tile3d

const (
	OP_NEGATIVE = "-"
	OP_POSITIVE = "+"
	OP_NOT      = "!"
)

const (
	FUNC_IS_NAN             = "isNaN"
	FUNC_IS_FINITE          = "isFinite"
	FUNC_IS_EXACTCLASS      = "isExactClass"
	FUNC_IS_CLASS           = "isClass"
	FUNC_GET_EXACTCLASSNAME = "getExactClassName"
	FUNC_BOOLEAN            = "Boolean"
	FUNC_NUMBER             = "Number"
	FUNC_STRING             = "String"
	FUNC_ABS                = "abs"
	FUNC_SQRT               = "sqrt"
	FUNC_COS                = "cos"
	FUNC_SIN                = "sin"
	FUNC_TAN                = "tan"
	FUNC_ACOS               = "acos"
	FUNC_ASIN               = "asin"
	FUNC_ATAN               = "atan"
	FUNC_RADIANS            = "radians"
	FUNC_DEGREES            = "degrees"
	FUNC_SIGN               = "sign"
	FUNC_FLOOR              = "floor"
	FUNC_CEIL               = "ceil"
	FUNC_ROUND              = "round"
	FUNC_EXP                = "exp"
	FUNC_EXP2               = "exp2"
	FUNC_LOG                = "log"
	FUNC_LOG2               = "log2"
	FUNC_FRACT              = "fract"
	FUNC_LENGTH             = "length"
	FUNC_NORMALIZE          = "normalize"
)

const (
	FUNC_ATAN2    = "atan2"
	FUNC_POW      = "pow"
	FUNC_MIN      = "min"
	FUNC_MAX      = "max"
	FUNC_DISTANCE = "distance"
	FUNC_DOT      = "dot"
	FUNC_CROSS    = "cross"
)

const (
	FUNC_CLAMP = "clamp"
	FUNC_MIX   = "mix"
)

type ExpressionType int32

const (
	Boolean   = ExpressionType(0)
	Null      = ExpressionType(1)
	Undefined = ExpressionType(2)
	Number    = ExpressionType(3)
	String    = ExpressionType(4)
	Array     = ExpressionType(5)
	vec2      = ExpressionType(6)
	vec3      = ExpressionType(7)
	vec4      = ExpressionType(8)
	RegExp    = ExpressionType(9)
)

type ExpressionNodeType int32

const (
	VARIABLE           = ExpressionNodeType(0)
	UNARY              = ExpressionNodeType(1)
	BINARY             = ExpressionNodeType(2)
	TERNARY            = ExpressionNodeType(3)
	CONDITIONAL        = ExpressionNodeType(4)
	MEMBER             = ExpressionNodeType(5)
	FUNCTION_CALL      = ExpressionNodeType(6)
	ARRAY              = ExpressionNodeType(7)
	REGEX              = ExpressionNodeType(8)
	VARIABLE_IN_STRING = ExpressionNodeType(9)
	LITERAL_NULL       = ExpressionNodeType(10)
	LITERAL_BOOLEAN    = ExpressionNodeType(11)
	LITERAL_NUMBER     = ExpressionNodeType(12)
	LITERAL_STRING     = ExpressionNodeType(13)
	LITERAL_COLOR      = ExpressionNodeType(14)
	LITERAL_VECTOR     = ExpressionNodeType(15)
	LITERAL_REGEX      = ExpressionNodeType(16)
	LITERAL_UNDEFINED  = ExpressionNodeType(17)
	BUILTIN_VARIABLE   = ExpressionNodeType(18)
)

const (
	OP_ADD        = "+"
	OP_SUB        = "-"
	OP_MUL        = "*"
	OP_DIV        = "/"
	OP_MOD        = "%"
	OP_EQ         = "==="
	OP_NEQ        = "!=="
	OP_GREATER    = ">"
	OP_GREATER_EQ = ">="
	OP_LESS       = "<"
	OP_LESS_EQ    = "<="
	OP_AND        = "&&"
	OP_OR         = "||"
	OP_REGEXP_NOT = "!~"
	OP_REGEXP     = "=~"
)

type Styling struct {
	Defines                  string `json:"defines"`
	Show                     string `json:"show"`
	Color                    string `json:"color"`
	PointSize                string `json:"pointSize,omitempty"`
	PointOutlineColor        string `json:"pointOutlineColor,omitempty"`
	PointOutlineWidth        string `json:"pointOutlineWidth,omitempty"`
	LabelColor               string `json:"labelColor,omitempty"`
	LabelOutlineColor        string `json:"labelOutlineColor,omitempty"`
	LabelOutlineWidth        string `json:"labelOutlineWidth,omitempty"`
	Font                     string `json:"font,omitempty"`
	LabelStyle               string `json:"labelStyle,omitempty"`
	BackgroundColor          string `json:"backgroundColor,omitempty"`
	BackgroundPadding        string `json:"backgroundPadding,omitempty"`
	BackgroundEnabled        string `json:"backgroundEnabled,omitempty"`
	ScaleByDistance          string `json:"scaleByDistance,omitempty"`
	TranslucencyByDistance   string `json:"translucencyByDistance,omitempty"`
	DistanceDisplayCondition string `json:"distanceDisplayCondition,omitempty"`
	HeightOffset             string `json:"heightOffset,omitempty"`
	AnchorLineEnabled        string `json:"anchorLineEnabled,omitempty"`
	AnchorLineColor          string `json:"anchorLineColor,omitempty"`
	Image                    string `json:"image,omitempty"`
	DisableDepthTestDistance string `json:"disableDepthTestDistance,omitempty"`
	HorizontalOrigin         string `json:"horizontalOrigin,omitempty"`
	VerticalOrigin           string `json:"verticalOrigin,omitempty"`
	LabelHorizontalOrigin    string `json:"labelHorizontalOrigin,omitempty"`
	LabelVerticalOrigin      string `json:"labelVerticalOrigin,omitempty"`
}

type Condition [2]string
type Conditions []Condition
type Defines map[string]string
type Meta map[string]string

type StyleExpression interface{}

type Expression struct {
}

type ConditionsExpression struct{}

type Operator interface{}

type UnaryOperators struct{}

type BinaryOperators struct{}

type Function interface{}

type UnaryFunction struct{}

type BinaryFunction struct{}

type TernaryFunctions struct{}
