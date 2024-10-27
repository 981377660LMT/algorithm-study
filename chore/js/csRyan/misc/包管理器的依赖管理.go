// 包管理器的依赖管理之如何解决依赖冲突和循环依赖的问题
// https://zhengtianbao.com/posts/package-manager/

// 1. 依赖关系描述
// 2. 依赖解析
// 3. 依赖冲突解决
// 出现了两个版本的 C，那应该保留哪个呢？根据 SemVer 规范的要求：
//
// 主版本号：不兼容的 API 修改
// 次版本号：向下兼容的功能性新增
// 修订号：向下兼容的问题修正
// !通过最小版本选择（Minimal Version Selection，MVS）算法来做选择(保留较大的版本号)。
// 主版本号冲突的情况，如果出现这种情况，那么就需要报错提示手动解决依赖了。
// 4. 循环依赖解决
// 检查当前依赖包是否在 visited map 中。如果存在，说明检测到循环依赖，直接跳过。
// 将当前依赖包添加到 visited map 中。
// 使用 defer 语句确保在函数退出前将当前依赖包从 visited map 中移除。

package main

import "sort"

func main() {

}

// Package represents a software package with a name, version, and dependencies.
type Package struct {
	Name         string
	Version      string
	Dependencies []*Package
}

// FlattenDependencies flattens the dependency tree into a map of unique dependencies.
func FlattenDependencies(pkg *Package) []*Package {
	flatDependencies := make(map[string]*Package)
	visited := make(map[string]struct{})
	collectDependencies(pkg, flatDependencies, visited)

	// Convert map to slice and sort
	var sortedDependencies []*Package
	for _, pkg := range flatDependencies {
		sortedDependencies = append(sortedDependencies, pkg)
	}
	// !将结果按字母排序后返回是为了方便测试，同样的 yarn 生成的 yarn.lock 文件按照字母排序是为了在依赖包变动时方便进行 git diff 比对差异。
	sort.Slice(sortedDependencies, func(i, j int) bool {
		return sortedDependencies[i].Name < sortedDependencies[j].Name
	})
	return sortedDependencies
}

// collectDependencies recursively collects dependencies and detects cycles.
func collectDependencies(pkg *Package, flatDependencies map[string]*Package, visited map[string]struct{}) {
	if _, has := visited[pkg.Name]; has {
		// detected cycle in dependencies, skip
		return
	}

	if existingPkg, exists := flatDependencies[pkg.Name]; exists {
		if compareVersions(pkg.Version, existingPkg.Version) > 0 {
			flatDependencies[pkg.Name] = pkg
		} else {
			return
		}
	} else {
		flatDependencies[pkg.Name] = pkg
	}

	visited[pkg.Name] = struct{}{}
	defer func() {
		delete(visited, pkg.Name)
	}()

	for _, dep := range pkg.Dependencies {
		collectDependencies(dep, flatDependencies, visited)
	}
}

// compareVersions compares two version strings.
func compareVersions(version1, version2 string) int {
	if version1 > version2 {
		return 1
	} else if version1 < version2 {
		return -1
	} else {
		return 0
	}
}
