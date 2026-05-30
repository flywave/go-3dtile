package next

import (
	"encoding/json"
	"testing"
)

func TestDefaultStyleJSON(t *testing.T) {
	s := DefaultStyle()
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Style
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.Show.BooleanExpression == nil {
		t.Error("Show.BooleanExpression should not be nil")
	}
	if decoded.Color.ColorExpression != "color('#FFFFFF')" {
		t.Errorf("Color.ColorExpression = %q, want color('#FFFFFF')", decoded.Color.ColorExpression)
	}
}

func TestStyleWithConditions(t *testing.T) {
	s := Style{
		Show: BooleanExpressionOrConditions{
			Conditions: &Conditions{
				Conditions: &[][2]string{
					{"${Height} > 10", "true"},
					{"true", "false"},
				},
			},
		},
		Color: ColorExpressionOrConditions{
			ColorExpression: "color('red')",
		},
		PointSize: NumberExpressionOrConditions{
			NumberExpression: json.RawMessage(`2.0`),
		},
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Style
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.Show.Conditions == nil {
		t.Fatal("Show.Conditions should not be nil")
	}
	if decoded.Show.BooleanExpression != nil {
		t.Error("Show.BooleanExpression should be nil when conditions are set")
	}
}

func TestBooleanExpressionJSON(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"true", `true`},
		{"false", `false`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b BooleanExpressionOrConditions
			if err := json.Unmarshal([]byte(tt.expr), &b); err != nil {
				t.Fatal("Unmarshal failed:", err)
			}
			if b.BooleanExpression == nil {
				t.Error("BooleanExpression should not be nil")
			}
			if b.Conditions != nil {
				t.Error("Conditions should be nil for simple expression")
			}

			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal("Marshal failed:", err)
			}
			if string(data) != tt.expr {
				t.Errorf("Marshal result = %s, want %s", string(data), tt.expr)
			}
		})
	}
}

func TestBooleanExpressionRoundTrip(t *testing.T) {
	expr := `"${Height} > 0"`
	var b BooleanExpressionOrConditions
	if err := json.Unmarshal([]byte(expr), &b); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}
	if b.BooleanExpression == nil {
		t.Error("BooleanExpression should not be nil")
	}

	data, err := json.Marshal(b)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded BooleanExpressionOrConditions
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Re-Unmarshal failed:", err)
	}
	if decoded.BooleanExpression == nil {
		t.Error("BooleanExpression should not be nil after round trip")
	}
}

func TestColorExpressionJSON(t *testing.T) {
	tests := []string{
		`"color('red')"`,
		`"rgb(255, 0, 0)"`,
	}
	for _, expr := range tests {
		t.Run(expr, func(t *testing.T) {
			var c ColorExpressionOrConditions
			if err := json.Unmarshal([]byte(expr), &c); err != nil {
				t.Fatal("Unmarshal failed:", err)
			}
			if c.ColorExpression == "" {
				t.Error("ColorExpression should not be empty")
			}
			if c.Conditions != nil {
				t.Error("Conditions should be nil")
			}

			data, err := json.Marshal(c)
			if err != nil {
				t.Fatal("Marshal failed:", err)
			}
			if string(data) != expr {
				t.Errorf("Marshal result = %s, want %s", string(data), expr)
			}
		})
	}
}

func TestNumberExpressionJSON(t *testing.T) {
	tests := []string{
		`1.0`,
		`2`,
		`"${Height} * 0.5"`,
	}
	for _, expr := range tests {
		t.Run(expr, func(t *testing.T) {
			var n NumberExpressionOrConditions
			if err := json.Unmarshal([]byte(expr), &n); err != nil {
				t.Fatal("Unmarshal failed:", err)
			}
			if n.NumberExpression == nil {
				t.Error("NumberExpression should not be nil")
			}
			if n.Conditions != nil {
				t.Error("Conditions should be nil")
			}

			data, err := json.Marshal(n)
			if err != nil {
				t.Fatal("Marshal failed:", err)
			}
			if string(data) != expr {
				t.Errorf("Marshal result = %s, want %s", string(data), expr)
			}
		})
	}
}

func TestStyleWithDefines(t *testing.T) {
	s := Style{
		Defines: map[string]string{
			"height": "${Height}",
		},
		Show: BooleanExpressionOrConditions{
			BooleanExpression: json.RawMessage(`true`),
		},
		Color: ColorExpressionOrConditions{
			ColorExpression: "color('blue')",
		},
		PointSize: NumberExpressionOrConditions{
			NumberExpression: json.RawMessage(`1.0`),
		},
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Style
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.Defines == nil || decoded.Defines["height"] != "${Height}" {
		t.Errorf("Defines = %v, want {height: ${Height}}", decoded.Defines)
	}
}

func TestStyleWithMeta(t *testing.T) {
	s := DefaultStyle()
	s.Meta = map[string]string{
		"name": "building",
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Style
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.Meta == nil || decoded.Meta["name"] != "building" {
		t.Errorf("Meta = %v, want {name: building}", decoded.Meta)
	}
}

func TestStyleWithExtensions(t *testing.T) {
	s := DefaultStyle()
	s.Extensions = map[string]json.RawMessage{
		"3DTILES_extension": json.RawMessage(`{"enabled": true}`),
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Style
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.Extensions == nil {
		t.Fatal("Extensions should not be nil")
	}
}

func TestConditionsUnmarshalInvalid(t *testing.T) {
	// invalid JSON should cause UnmarshalJSON to fail
	var c ColorExpressionOrConditions
	err := json.Unmarshal([]byte(`{"invalid": true}`), &c)
	if err == nil {
		t.Error("Expected error for invalid color expression, got nil")
	}
}

func TestStyleDefinesMarshalUnmarshal(t *testing.T) {
	// Verify that empty defines are omitted
	s := DefaultStyle()
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if _, exists := raw["defines"]; exists {
		t.Error("defines should be omitted when empty")
	}
}
