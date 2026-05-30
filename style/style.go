// Package style implements the 3D Tiles styling specification.
// It provides a unified type system shared by both the legacy (tile3d) and
// next-generation (next) API surfaces, supporting both simple string expressions
// and structured condition-based styles.
package style

import "encoding/json"

// Style defines how features are displayed using 3D Tiles styling expressions.
// Each expression field accepts either a simple expression string (e.g. "1.0",
// "true", "color('red')") or a conditions array [[cond, val], ...].
type Style struct {
	Defines                  map[string]string                 `json:"defines,omitempty"`
	Show                     BooleanExpressionOrConditions     `json:"show"`
	Color                    ColorExpressionOrConditions       `json:"color"`
	Meta                     map[string]string                 `json:"meta,omitempty"`
	PointSize                NumberExpressionOrConditions      `json:"pointSize"`
	PointOutlineColor        *StringExpressionOrConditions     `json:"pointOutlineColor,omitempty"`
	PointOutlineWidth        *StringExpressionOrConditions     `json:"pointOutlineWidth,omitempty"`
	LabelColor               *StringExpressionOrConditions     `json:"labelColor,omitempty"`
	LabelOutlineColor        *StringExpressionOrConditions     `json:"labelOutlineColor,omitempty"`
	LabelOutlineWidth        *StringExpressionOrConditions     `json:"labelOutlineWidth,omitempty"`
	Font                     *StringExpressionOrConditions     `json:"font,omitempty"`
	LabelStyle               *StringExpressionOrConditions     `json:"labelStyle,omitempty"`
	BackgroundColor          *StringExpressionOrConditions     `json:"backgroundColor,omitempty"`
	BackgroundPadding        *StringExpressionOrConditions     `json:"backgroundPadding,omitempty"`
	BackgroundEnabled        *StringExpressionOrConditions     `json:"backgroundEnabled,omitempty"`
	ScaleByDistance          *StringExpressionOrConditions     `json:"scaleByDistance,omitempty"`
	TranslucencyByDistance   *StringExpressionOrConditions     `json:"translucencyByDistance,omitempty"`
	DistanceDisplayCondition *StringExpressionOrConditions     `json:"distanceDisplayCondition,omitempty"`
	HeightOffset             *StringExpressionOrConditions     `json:"heightOffset,omitempty"`
	AnchorLineEnabled        *StringExpressionOrConditions     `json:"anchorLineEnabled,omitempty"`
	AnchorLineColor          *StringExpressionOrConditions     `json:"anchorLineColor,omitempty"`
	Image                    *StringExpressionOrConditions     `json:"image,omitempty"`
	DisableDepthTestDistance *StringExpressionOrConditions     `json:"disableDepthTestDistance,omitempty"`
	HorizontalOrigin         *StringExpressionOrConditions     `json:"horizontalOrigin,omitempty"`
	VerticalOrigin           *StringExpressionOrConditions     `json:"verticalOrigin,omitempty"`
	LabelHorizontalOrigin    *StringExpressionOrConditions     `json:"labelHorizontalOrigin,omitempty"`
	LabelVerticalOrigin      *StringExpressionOrConditions     `json:"labelVerticalOrigin,omitempty"`
	Extensions               map[string]json.RawMessage        `json:"extensions,omitempty"`
	Extra                    json.RawMessage                   `json:"extra,omitempty"`
}

// DefaultStyle returns a new Style with default values.
func DefaultStyle() Style {
	return Style{
		Show: BooleanExpressionOrConditions{
			BooleanExpression: json.RawMessage(`true`),
		},
		Color: ColorExpressionOrConditions{
			ColorExpression: "color('#FFFFFF')",
		},
		PointSize: NumberExpressionOrConditions{
			NumberExpression: json.RawMessage(`1.0`),
		},
	}
}

// Expression is a 3D Tiles style expression string.
type Expression = string

// BooleanExpression represents a boolean or string expression (json.RawMessage).
type BooleanExpression = json.RawMessage

// ColorExpression represents a color expression string.
type ColorExpression = string

// NumberExpression represents a number expression (json.RawMessage).
type NumberExpression = json.RawMessage
