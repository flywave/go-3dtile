package subtree

import (
	"testing"
)

func TestSubtreeCreator_GenerateSubtreefile(t *testing.T) {
	// 创建测试用的tiles
	tiles := []Tile{
		{Z: 0, X: 0, Y: 0, Available: true},
		{Z: 1, X: 0, Y: 0, Available: true},
		{Z: 1, X: 1, Y: 1, Available: true},
	}

	// 生成子树文件
	subtreeBytes := GenerateSubtreefile(tiles)

	// 验证结果不为空
	if len(subtreeBytes) == 0 {
		t.Error("Expected non-empty subtree bytes")
	}
}

func TestSubtreeCreator_GenerateSubtreefiles(t *testing.T) {
	// 创建测试用的tiles
	tiles := []Tile{
		{Z: 0, X: 0, Y: 0, Available: true},
		{Z: 1, X: 0, Y: 0, Available: true},
		{Z: 1, X: 1, Y: 1, Available: true},
	}

	// 生成子树文件
	subtreeFiles := GenerateSubtreefiles(tiles)

	// 验证结果不为空
	if len(subtreeFiles) == 0 {
		t.Error("Expected non-empty subtree files")
	}

	// 验证根节点存在
	rootTile := Tile{Z: 0, X: 0, Y: 0}
	if _, exists := subtreeFiles[rootTile]; !exists {
		t.Error("Expected root tile in subtree files")
	}
}

func TestSubtreeCreator_Fill(t *testing.T) {
	// 测试填充功能
	availability := "101"
	targetLength := 5
	result := Fill(availability, targetLength)

	if len(result) != targetLength {
		t.Errorf("Expected length %d, got %d", targetLength, len(result))
	}

	if result != "10100" {
		t.Errorf("Expected '10100', got '%s'", result)
	}

	// 测试不需要填充的情况
	availability = "10101"
	targetLength = 3
	result = Fill(availability, targetLength)

	if result != "10101" {
		t.Errorf("Expected '10101', got '%s'", result)
	}
}

func TestSubtreeCreator_GetSubtreeTiles(t *testing.T) {
	// 创建测试用的tiles
	tiles := []Tile{
		{Z: 0, X: 0, Y: 0, Available: true},
		{Z: 1, X: 0, Y: 0, Available: true},
		{Z: 1, X: 1, Y: 1, Available: true},
		{Z: 2, X: 2, Y: 2, Available: true},
	}

	// 获取子树tiles
	tile := Tile{Z: 1, X: 0, Y: 0}
	subtreeTiles := GetSubtreeTiles(tiles, tile)

	// 验证结果不为空
	if len(subtreeTiles) == 0 {
		t.Error("Expected non-empty subtree tiles")
	}
}

func TestSubtreeCreator_GetRelativeTile(t *testing.T) {
	// 测试相对位置计算
	from := Tile{Z: 1, X: 1, Y: 1}
	to := Tile{Z: 2, X: 2, Y: 2}

	relative := GetRelativeTile(from, to)

	// 验证相对位置
	expectedX := 0 // 2 - (1 * 2) = 0
	expectedY := 0 // 2 - (1 * 2) = 0
	expectedZ := 1 // 2 - 1 = 1

	if relative.X != expectedX || relative.Y != expectedY || relative.Z != expectedZ {
		t.Errorf("Expected relative tile (%d, %d, %d), got (%d, %d, %d)",
			expectedX, expectedY, expectedZ, relative.X, relative.Y, relative.Z)
	}
}

func TestTile_HasChild(t *testing.T) {
	// 测试父级tile是否有子tile
	parent := Tile{Z: 0, X: 0, Y: 0}
	child := Tile{Z: 1, X: 0, Y: 0}

	if !parent.HasChild(child) {
		t.Error("Expected parent to have child")
	}

	// 测试非子节点
	notChild := Tile{Z: 1, X: 2, Y: 2}
	if parent.HasChild(notChild) {
		t.Error("Expected parent to not have child")
	}

	// 测试同级节点
	sameLevel := Tile{Z: 0, X: 1, Y: 1}
	if parent.HasChild(sameLevel) {
		t.Error("Expected parent to not have same level tile as child")
	}

	// 测试父级节点
	parentTile := Tile{Z: 1, X: 0, Y: 0}
	childTile := Tile{Z: 0, X: 0, Y: 0}
	if parentTile.HasChild(childTile) {
		t.Error("Expected parent to not have child with lower Z level")
	}
}

func TestGetMortonIndices(t *testing.T) {
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
}

func TestGetMortonIndices_Empty(t *testing.T) {
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

func TestGetMortonIndexAsBytes(t *testing.T) {
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

func TestGetMortonIndexAsBytes_Empty(t *testing.T) {
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

func TestGenerateSubtreefile(t *testing.T) {
	// 创建测试用的tiles
	tiles := []Tile{
		{Z: 0, X: 0, Y: 0, Available: true},
		{Z: 1, X: 0, Y: 0, Available: true},
		{Z: 1, X: 1, Y: 1, Available: true},
	}

	// 生成子树文件
	subtreeBytes := GenerateSubtreefile(tiles)

	// 验证结果不为空
	if len(subtreeBytes) == 0 {
		t.Error("Expected non-empty subtree bytes")
	}
}

func TestGenerateSubtreefile_Empty(t *testing.T) {
	// 创建空的tiles
	tiles := []Tile{}

	// 生成子树文件
	subtreeBytes := GenerateSubtreefile(tiles)

	// 验证结果不为空（即使没有tiles也应该生成有效的子树文件）
	if len(subtreeBytes) == 0 {
		t.Error("Expected non-empty subtree bytes even for empty tiles")
	}
}

func TestGenerateSubtreefiles(t *testing.T) {
	// 创建测试用的tiles
	tiles := []Tile{
		{Z: 0, X: 0, Y: 0, Available: true},
		{Z: 1, X: 0, Y: 0, Available: true},
		{Z: 1, X: 1, Y: 1, Available: true},
		{Z: 2, X: 0, Y: 0, Available: true},
	}

	// 生成子树文件
	subtreeFiles := GenerateSubtreefiles(tiles)

	// 验证结果不为空
	if len(subtreeFiles) == 0 {
		t.Error("Expected non-empty subtree files")
	}

	// 验证根节点存在
	rootTile := Tile{Z: 0, X: 0, Y: 0}
	if _, exists := subtreeFiles[rootTile]; !exists {
		t.Error("Expected root tile in subtree files")
	}
}

func TestGenerateSubtreefiles_SingleLevel(t *testing.T) {
	// 创建单级的tiles
	tiles := []Tile{
		{Z: 0, X: 0, Y: 0, Available: true},
	}

	// 生成子树文件
	subtreeFiles := GenerateSubtreefiles(tiles)

	// 验证结果不为空
	if len(subtreeFiles) == 0 {
		t.Error("Expected non-empty subtree files")
	}

	// 验证根节点存在
	rootTile := Tile{Z: 0, X: 0, Y: 0}
	if _, exists := subtreeFiles[rootTile]; !exists {
		t.Error("Expected root tile in subtree files")
	}
}

func TestFill(t *testing.T) {
	// 测试填充功能
	availability := "101"
	targetLength := 5
	result := Fill(availability, targetLength)

	if len(result) != targetLength {
		t.Errorf("Expected length %d, got %d", targetLength, len(result))
	}

	if result != "10100" {
		t.Errorf("Expected '10100', got '%s'", result)
	}

	// 测试不需要填充的情况
	availability = "10101"
	targetLength = 3
	result = Fill(availability, targetLength)

	if result != "10101" {
		t.Errorf("Expected '10101', got '%s'", result)
	}

	// 测试空字符串
	availability = ""
	targetLength = 3
	result = Fill(availability, targetLength)

	if result != "000" {
		t.Errorf("Expected '000', got '%s'", result)
	}
}

func TestGetSubtreeTiles(t *testing.T) {
	// 创建测试用的tiles
	tiles := []Tile{
		{Z: 0, X: 0, Y: 0, Available: true},
		{Z: 1, X: 0, Y: 0, Available: true},
		{Z: 1, X: 1, Y: 1, Available: true},
		{Z: 2, X: 2, Y: 2, Available: true},
	}

	// 获取子树tiles
	tile := Tile{Z: 1, X: 0, Y: 0}
	subtreeTiles := GetSubtreeTiles(tiles, tile)

	// 验证结果不为空
	if len(subtreeTiles) == 0 {
		t.Error("Expected non-empty subtree tiles")
	}

	// 验证根tile存在
	if len(subtreeTiles) < 1 {
		t.Error("Expected at least root tile in subtree tiles")
	}
}

func TestGetRelativeTile(t *testing.T) {
	// 测试相对位置计算
	from := Tile{Z: 1, X: 1, Y: 1}
	to := Tile{Z: 2, X: 2, Y: 2}

	relative := GetRelativeTile(from, to)

	// 验证相对位置
	expectedX := 0 // 2 - (1 * 2) = 0
	expectedY := 0 // 2 - (1 * 2) = 0
	expectedZ := 1 // 2 - 1 = 1

	if relative.X != expectedX || relative.Y != expectedY || relative.Z != expectedZ {
		t.Errorf("Expected relative tile (%d, %d, %d), got (%d, %d, %d)",
			expectedX, expectedY, expectedZ, relative.X, relative.Y, relative.Z)
	}

	// 测试相同级别
	from = Tile{Z: 0, X: 0, Y: 0}
	to = Tile{Z: 0, X: 1, Y: 1}

	relative = GetRelativeTile(from, to)

	// 验证相对位置
	expectedX = 1 // 1 - (0 * 1) = 1
	expectedY = 1 // 1 - (0 * 1) = 1
	expectedZ = 0 // 0 - 0 = 0

	if relative.X != expectedX || relative.Y != expectedY || relative.Z != expectedZ {
		t.Errorf("Expected relative tile (%d, %d, %d), got (%d, %d, %d)",
			expectedX, expectedY, expectedZ, relative.X, relative.Y, relative.Z)
	}
}
