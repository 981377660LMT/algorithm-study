// 脱离LSM，抽象出 Size-Tiered Compaction（大小分层合并）  和 Leveled Compaction（分层合并） 的算法实现
//
// Size-Tiered Compaction（STC）：
// 将多个相同或相近大小的数据文件合并成一个更大的文件，适合高写入频率的场景。
// Leveled Compaction（LC）：
// 类似二进制分组。
// 将数据文件分配到不同的层级中，每个层级中的文件互不重叠，适合需要稳定读取性能的场景。

package main

import "fmt"

// sizeTieredCompaction 执行 Size-Tiered Compaction
//  返回值为合并后的文件大小切片
//  files: 输入的文件大小切片
//  mergeThreshold: 合并阈值，当某一组的文件数量 >= mergeThreshold 时进行合并
func sizeTieredCompaction(files []int, mergeThreshold int) []int {
	// 1. 按大小分组
	sizeGroups := groupBySize(files)

	// 2. 记录是否进行了合并，以决定是否需要继续
	merged := false

	// 3. 遍历每个大小组，检查是否需要合并
	for size, group := range sizeGroups {
		if len(group) >= mergeThreshold {
			// 进行合并
			mergedSize := mergeFiles(group)
			// 添加新的合并后的文件
			files = append(files, mergedSize)
			// 从原始文件列表中删除被合并的文件
			files = removeFilesBySizes(files, group)
			// 标记已进行合并
			merged = true
			fmt.Printf("Merged files of size %d: %v into new file of size %d\n", size, group, mergedSize)
		}
	}

	// 4. 如果进行了合并，递归调用以处理可能的新合并机会
	if merged {
		return sizeTieredCompaction(files, mergeThreshold)
	}

	// 5. 返回最终合并后的文件列表
	return files
}

// groupBySize 将文件按大小分组，返回大小到文件列表的映射
func groupBySize(files []int) map[int][]int {
	sizeGroups := make(map[int][]int)
	for _, size := range files {
		sizeGroups[size] = append(sizeGroups[size], size)
	}
	return sizeGroups
}

// mergeFiles 合并一组文件，返回合并后的文件大小
// 这里的合并逻辑是简单的求和，可以根据需求调整
func mergeFiles(group []int) int {
	total := 0
	for _, size := range group {
		total += size
	}
	return total
}

// removeFilesBySizes 从文件列表中移除指定大小的文件
func removeFilesBySizes(files []int, sizes []int) []int {
	toRemove := make(map[int]int) // size -> count
	for _, size := range sizes {
		toRemove[size]++
	}

	result := []int{}
	for _, size := range files {
		if count, exists := toRemove[size]; exists && count > 0 {
			toRemove[size]--
			continue
		}
		result = append(result, size)
	}
	return result
}

// leveledCompaction 执行 Leveled Compaction
// files: 输入的文件大小切片
// maxFilesPerLevel: 每层级最大文件数
// sizeMultiplier: 层级大小倍数（例如，Level 1 的大小限制是 Level 0 的 sizeMultiplier 倍）
// 返回值为合并后的文件大小切片
func leveledCompaction(files []int, maxFilesPerLevel int, sizeMultiplier int) []int {
	// 1. 初始化层级映射，level -> 文件大小列表
	levels := make(map[int][]int)
	for _, size := range files {
		// 初始分配到 Level 0
		levels[0] = append(levels[0], size)
	}

	level := 0
	for {
		// 获取当前层级的所有文件
		currentLevelFiles := levels[level]
		if len(currentLevelFiles) > maxFilesPerLevel {
			// 需要合并到下一层级
			nextLevel := level + 1
			// 合并当前层级的所有文件
			mergedSize := mergeAndDedup(currentLevelFiles)
			// 添加合并后的文件到下一层级
			levels[nextLevel] = append(levels[nextLevel], mergedSize)
			fmt.Printf("Merged %d files from Level %d into Level %d as size %d\n", len(currentLevelFiles), level, nextLevel, mergedSize)
			// 清空当前层级的文件
			levels[level] = []int{}
			// 继续检查当前层级（可能仍有超过限制的文件）
			continue
		} else {
			// 检查是否需要继续向上层级合并
			// 计算下一层级的大小限制
			expectedSize := calculateLevelSize(level+1, sizeMultiplier)
			// 计算下一层级的总大小
			totalSize := sum(levels[level+1])
			if level+1 > 0 && totalSize > expectedSize {
				level++
				continue
			}
		}
		// 没有需要合并的层级，退出循环
		break
	}

	// 2. 收集所有层级的文件大小
	mergedFiles := []int{}
	for _, levelFiles := range levels {
		mergedFiles = append(mergedFiles, levelFiles...)
	}

	// 3. 返回合并后的文件列表
	return mergedFiles
}

// mergeAndDedup 合并并去重文件大小
// 这里的合并逻辑是简单的求和，可以根据需求调整
func mergeAndDedup(files []int) int {
	// 假设合并所有文件为一个文件，大小为它们的总和
	total := 0
	for _, size := range files {
		total += size
	}
	return total
}

// calculateLevelSize 计算指定层级的大小限制
func calculateLevelSize(level int, multiplier int) int {
	if level == 0 {
		return 0 // Level 0 没有大小限制
	}
	return pow(multiplier, level-1) * 10 // 示例：每层级大小倍数为10
}

// pow 计算整数的幂
func pow(a, b int) int {
	result := 1
	for i := 0; i < b; i++ {
		result *= a
	}
	return result
}

// sum 计算切片中所有整数的和
func sum(nums []int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}

// leveledCompaction 示例使用
func leveledCompactionExample() {
	fmt.Println("\n=== Leveled Compaction Example ===")
	// 初始文件大小列表
	files := []int{5, 5, 5, 10, 10, 15}
	// 每层级最大文件数为2，层级大小倍数为10
	maxFilesPerLevel := 2
	sizeMultiplier := 10

	fmt.Println("Before Leveled Compaction:", files)
	mergedFiles := leveledCompaction(files, maxFilesPerLevel, sizeMultiplier)
	fmt.Println("After Leveled Compaction:", mergedFiles)
}
