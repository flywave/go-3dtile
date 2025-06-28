package next

import (
	"encoding/json"
)

// Style represents a 3D Tiles style
type Style struct {
	Defines    map[string]string             `json:"defines,omitempty"`
	Show       BooleanExpressionOrConditions `json:"show"`
	Color      ColorExpressionOrConditions   `json:"color"`
	Meta       map[string]string             `json:"meta,omitempty"`
	PointSize  NumberExpressionOrConditions  `json:"pointSize"`
	Extensions map[string]json.RawMessage    `json:"extensions,omitempty"`
	Extra      json.RawMessage               `json:"extra,omitempty"`
}

// DefaultStyle returns a new Style with default values
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

// Expression represents a valid 3D Tiles style expression
type Expression = string

// BooleanExpression represents a boolean or string expression
type BooleanExpression = json.RawMessage

// ColorExpression represents a color expression string
type ColorExpression = string

// NumberExpression represents a number expression
type NumberExpression = json.RawMessage

// Meta represents feature metadata properties
type Meta = map[string]string

// BooleanExpressionOrConditions represents either a boolean expression or conditions
type BooleanExpressionOrConditions struct {
	BooleanExpression BooleanExpression `json:"-"`
	Conditions        *Conditions       `json:"-"`
}

// Implement custom JSON marshaling
func (b BooleanExpressionOrConditions) MarshalJSON() ([]byte, error) {
	if b.Conditions != nil {
		return json.Marshal(b.Conditions)
	}
	return json.Marshal(b.BooleanExpression)
}

// Implement custom JSON unmarshaling
func (b *BooleanExpressionOrConditions) UnmarshalJSON(data []byte) error {
	// First try to unmarshal as Conditions
	var conditions Conditions
	if err := json.Unmarshal(data, &conditions); err == nil && conditions.Conditions != nil {
		b.Conditions = &conditions
		return nil
	}

	// Otherwise unmarshal as BooleanExpression
	var expr BooleanExpression
	if err := json.Unmarshal(data, &expr); err != nil {
		return err
	}
	b.BooleanExpression = expr
	return nil
}

// NumberExpressionOrConditions represents either a number expression or conditions
type NumberExpressionOrConditions struct {
	NumberExpression NumberExpression `json:"-"`
	Conditions       *Conditions      `json:"-"`
}

// Implement custom JSON marshaling
func (n NumberExpressionOrConditions) MarshalJSON() ([]byte, error) {
	if n.Conditions != nil {
		return json.Marshal(n.Conditions)
	}
	return json.Marshal(n.NumberExpression)
}

// Implement custom JSON unmarshaling
func (n *NumberExpressionOrConditions) UnmarshalJSON(data []byte) error {
	// First try to unmarshal as Conditions
	var conditions Conditions
	if err := json.Unmarshal(data, &conditions); err == nil && conditions.Conditions != nil {
		n.Conditions = &conditions
		return nil
	}

	// Otherwise unmarshal as NumberExpression
	var expr NumberExpression
	if err := json.Unmarshal(data, &expr); err != nil {
		return err
	}
	n.NumberExpression = expr
	return nil
}

// ColorExpressionOrConditions represents either a color expression or conditions
type ColorExpressionOrConditions struct {
	ColorExpression ColorExpression `json:"-"`
	Conditions      *Conditions     `json:"-"`
}

// Implement custom JSON marshaling
func (c ColorExpressionOrConditions) MarshalJSON() ([]byte, error) {
	if c.Conditions != nil {
		return json.Marshal(c.Conditions)
	}
	return json.Marshal(c.ColorExpression)
}

// Implement custom JSON unmarshaling
func (c *ColorExpressionOrConditions) UnmarshalJSON(data []byte) error {
	// First try to unmarshal as Conditions
	var conditions Conditions
	if err := json.Unmarshal(data, &conditions); err == nil && conditions.Conditions != nil {
		c.Conditions = &conditions
		return nil
	}

	// Otherwise unmarshal as ColorExpression
	var expr ColorExpression
	if err := json.Unmarshal(data, &expr); err != nil {
		return err
	}
	c.ColorExpression = expr
	return nil
}

// Conditions represents a series of conditional expressions
type Conditions struct {
	Conditions *[][2]string `json:"conditions,omitempty"`
}
