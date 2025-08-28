package subtree

import (
	"testing"
)

func TestMortonIndex_GetMortonIndices(t *testing.T) {
	// 创建测试用的tiles
	tiles := []Tile{
		{Z: 0, X: 0, Y: 0, Available: true},
		{Z: 1, X: 0, Y: 0, Available: true},
		{Z: 1, X: 1, Y: 1, Available: true},
	}

	// 计算Morton索引
	tileAvailability, contentAvailability := GetMortonIndices(tiles)

	// 验证结果不为空
	if tileAvailability == "" {
		t.Error("Expected non-empty tile availability")
	}

	if contentAvailability == "" {
		t.Error("Expected non-empty content availability")
	}

	// 验证返回值类型
	if len(tileAvailability) == 0 {
		t.Error("Expected tile availability string")
	}

	if len(contentAvailability) == 0 {
		t.Error("Expected content availability string")
	}
}

func TestMortonIndex_GetMortonIndexAsBytes(t *testing.T) {
	// 创建测试用的tiles
	tiles := []Tile{
		{Z: 0, X: 0, Y: 0, Available: true},
		{Z: 1, X: 0, Y: 0, Available: true},
		{Z: 1, X: 1, Y: 1, Available: true},
	}

	// 计算Morton索引为字节数组
	tileAvailability, contentAvailability := GetMortonIndexAsBytes(tiles)

	// 验证结果不为空
	if len(tileAvailability) == 0 {
		t.Error("Expected non-empty tile availability byte array")
	}

	if len(contentAvailability) == 0 {
		t.Error("Expected non-empty content availability byte array")
	}
}

func TestMortonIndex_GetMortonIndices_Empty(t *testing.T) {
	// 创建空的tiles
	tiles := []Tile{}

	// 计算Morton索引
	tileAvailability, contentAvailability := GetMortonIndices(tiles)

	// 验证结果为空但不为nil
	if tileAvailability != "" {
		t.Errorf("Expected empty tile availability, got %s", tileAvailability)
	}

	if contentAvailability != "" {
		t.Errorf("Expected empty content availability, got %s", contentAvailability)
	}
}

func TestMortonIndex_GetMortonIndexAsBytes_Empty(t *testing.T) {
	// 创建空的tiles
	tiles := []Tile{}

	// 计算Morton索引为字节数组
	tileAvailability, contentAvailability := GetMortonIndexAsBytes(tiles)

	// 验证结果为空数组
	if len(tileAvailability) != 0 {
		t.Errorf("Expected empty tile availability byte array, got length %d", len(tileAvailability))
	}

	if len(contentAvailability) != 0 {
		t.Errorf("Expected empty content availability byte array, got length %d", len(contentAvailability))
	}
}
