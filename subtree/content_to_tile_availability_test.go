package subtree

import (
	"testing"
)

func TestContentToTileAvailability_GetTileAvailabilityLevels(t *testing.T) {
	// 创建测试用的内容可用性级别
	contentLevels := make(AvailabilityLevels, 0)

	// 添加级别0
	level0 := NewAvailabilityLevel(0, Quadtree)
	contentLevels = append(contentLevels, level0)

	// 添加级别1
	level1 := NewAvailabilityLevel(1, Quadtree)
	contentLevels = append(contentLevels, level1)

	// 在级别1的内容中设置一些值
	if level1.BitArray2D != nil {
		level1.BitArray2D.Set(0, 0, true)
		level1.BitArray2D.Set(1, 1, true)
	}

	// 计算瓦片可用性级别
	tileLevels := GetTileAvailabilityLevels(contentLevels)

	// 验证结果
	if len(tileLevels) != 2 {
		t.Errorf("Expected 2 tile levels, got %d", len(tileLevels))
	}

	// 验证级别0存在
	tileLevel0 := tileLevels.GetLevel(0)
	if tileLevel0 == nil {
		t.Error("Expected level 0 in tile levels")
	}

	// 验证级别1存在
	tileLevel1 := tileLevels.GetLevel(1)
	if tileLevel1 == nil {
		t.Error("Expected level 1 in tile levels")
	}

	// 验证级别1的瓦片设置
	if tileLevel1 != nil && tileLevel1.BitArray2D != nil {
		if !tileLevel1.BitArray2D.Get(0, 0) {
			t.Error("Expected (0,0) to be true in level 1 tile availability")
		}
		if !tileLevel1.BitArray2D.Get(1, 1) {
			t.Error("Expected (1,1) to be true in level 1 tile availability")
		}
	}

	// 验证父级瓦片设置 (0,0) >> 1 = 0
	if tileLevel0 != nil && tileLevel0.BitArray2D != nil {
		if !tileLevel0.BitArray2D.Get(0, 0) {
			t.Error("Expected (0,0) to be true in level 0 tile availability (parent of content)")
		}
	}
}

func TestContentToTileAvailability_GetTileAvailabilityLevels_SingleLevel(t *testing.T) {
	// 创建只有级别0的测试用例
	contentLevels := make(AvailabilityLevels, 0)
	level0 := NewAvailabilityLevel(0, Quadtree)

	// 在级别0的内容中设置值
	if level0.BitArray2D != nil {
		level0.BitArray2D.Set(0, 0, true)
	}
	contentLevels = append(contentLevels, level0)

	// 计算瓦片可用性级别
	tileLevels := GetTileAvailabilityLevels(contentLevels)

	// 验证结果
	if len(tileLevels) != 1 {
		t.Errorf("Expected 1 tile level, got %d", len(tileLevels))
	}

	// 验证级别0的瓦片设置
	tileLevel0 := tileLevels.GetLevel(0)
	if tileLevel0 != nil && tileLevel0.BitArray2D != nil {
		if !tileLevel0.BitArray2D.Get(0, 0) {
			t.Error("Expected (0,0) to be true in level 0 tile availability")
		}
	}
}

func TestContentToTileAvailability_GetTileAvailabilityLevels_Empty(t *testing.T) {
	// 创建空的内容可用性级别
	contentLevels := make(AvailabilityLevels, 0)

	// 计算瓦片可用性级别
	tileLevels := GetTileAvailabilityLevels(contentLevels)

	// 验证结果
	if len(tileLevels) != 0 {
		t.Errorf("Expected 0 tile levels, got %d", len(tileLevels))
	}
}
