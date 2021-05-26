package tile3d

const (
	OP_NEGATIVE = "-"
	OP_POSITIVE = "+"
	OP_NOT      = "!"
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
	VAR_BOOLEAN   = ExpressionType(0)
	VAR_NULL      = ExpressionType(1)
	VAR_UNDEFINED = ExpressionType(2)
	VAR_NUMBER    = ExpressionType(3)
	VAR_STRING    = ExpressionType(4)
	VAR_ARRAY     = ExpressionType(5)
	VAR_VEC2      = ExpressionType(6)
	VAR_VEC3      = ExpressionType(7)
	VAR_VEC4      = ExpressionType(8)
	VAR_REGEXP    = ExpressionType(9)
)

type ExpressionNodeType int32

const (
	EXP_VARIABLE           = ExpressionNodeType(0)
	EXP_UNARY              = ExpressionNodeType(1)
	EXP_BINARY             = ExpressionNodeType(2)
	EXP_TERNARY            = ExpressionNodeType(3)
	EXP_CONDITIONAL        = ExpressionNodeType(4)
	EXP_MEMBER             = ExpressionNodeType(5)
	EXP_FUNCTION_CALL      = ExpressionNodeType(6)
	EXP_ARRAY              = ExpressionNodeType(7)
	EXP_REGEX              = ExpressionNodeType(8)
	EXP_VARIABLE_IN_STRING = ExpressionNodeType(9)
	EXP_LITERAL_NULL       = ExpressionNodeType(10)
	EXP_LITERAL_BOOLEAN    = ExpressionNodeType(11)
	EXP_LITERAL_NUMBER     = ExpressionNodeType(12)
	EXP_LITERAL_STRING     = ExpressionNodeType(13)
	EXP_LITERAL_COLOR      = ExpressionNodeType(14)
	EXP_LITERAL_VECTOR     = ExpressionNodeType(15)
	EXP_LITERAL_REGEX      = ExpressionNodeType(16)
	EXP_LITERAL_UNDEFINED  = ExpressionNodeType(17)
	EXP_BUILTIN_VARIABLE   = ExpressionNodeType(18)
)

type Expression string

type Styling struct {
	Defines                  Expression `json:"defines"`
	Show                     Expression `json:"show"`
	Color                    Expression `json:"color"`
	PointSize                Expression `json:"pointSize,omitempty"`
	PointOutlineColor        Expression `json:"pointOutlineColor,omitempty"`
	PointOutlineWidth        Expression `json:"pointOutlineWidth,omitempty"`
	LabelColor               Expression `json:"labelColor,omitempty"`
	LabelOutlineColor        Expression `json:"labelOutlineColor,omitempty"`
	LabelOutlineWidth        Expression `json:"labelOutlineWidth,omitempty"`
	Font                     Expression `json:"font,omitempty"`
	LabelStyle               Expression `json:"labelStyle,omitempty"`
	BackgroundColor          Expression `json:"backgroundColor,omitempty"`
	BackgroundPadding        Expression `json:"backgroundPadding,omitempty"`
	BackgroundEnabled        Expression `json:"backgroundEnabled,omitempty"`
	ScaleByDistance          Expression `json:"scaleByDistance,omitempty"`
	TranslucencyByDistance   Expression `json:"translucencyByDistance,omitempty"`
	DistanceDisplayCondition Expression `json:"distanceDisplayCondition,omitempty"`
	HeightOffset             Expression `json:"heightOffset,omitempty"`
	AnchorLineEnabled        Expression `json:"anchorLineEnabled,omitempty"`
	AnchorLineColor          Expression `json:"anchorLineColor,omitempty"`
	Image                    Expression `json:"image,omitempty"`
	DisableDepthTestDistance Expression `json:"disableDepthTestDistance,omitempty"`
	HorizontalOrigin         Expression `json:"horizontalOrigin,omitempty"`
	VerticalOrigin           Expression `json:"verticalOrigin,omitempty"`
	LabelHorizontalOrigin    Expression `json:"labelHorizontalOrigin,omitempty"`
	LabelVerticalOrigin      Expression `json:"labelVerticalOrigin,omitempty"`
}
