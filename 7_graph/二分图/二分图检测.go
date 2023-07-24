// 二分图检测

package main

func IsBipartite(n int, graph [][]int) (colors []int, ok bool) {
	colors = make([]int, n)
	for i := range colors {
		colors[i] = -1
	}
	bfs := func(start int) bool {
		colors[start] = 0
		queue := []int{start}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for _, next := range graph[cur] {
				if colors[next] == -1 {
					colors[next] = colors[cur] ^ 1
					queue = append(queue, next)
				} else if colors[next] == colors[cur] {
					return false
				}
			}
		}
		return true
	}

	for i := range colors {
		if colors[i] == -1 && !bfs(i) {
			return nil, false
		}
	}
	return colors, true
}

// def isBipartite3(n: int, adjList: List[List[int]]) -> bool:
//     """二分图检测 扩展域并查集"""
//     OFFSET = int(1e9)
//     uf = UnionFind()
//     for cur in range(n):
//         for next in adjList[cur]:
//             uf.union(cur, next + OFFSET)
//             uf.union(cur + OFFSET, next)
//             if uf.isConnected(cur, next):
//                 return False
//     return True

// T = TypeVar("T")

// class UnionFind(Generic[T]):
//     def __init__(self, iterable: Optional[Iterable[T]] = None):
//         self.count = 0
//         self.parent = dict()
//         self.rank = defaultdict(lambda: 1)
//         for item in iterable or []:
//             self.add(item)

//     def union(self, key1: T, key2: T) -> bool:
//         """rank一样时 默认key2作为key1的父节点"""
//         root1 = self.find(key1)
//         root2 = self.find(key2)
//         if root1 == root2:
//             return False
//         if self.rank[root1] > self.rank[root2]:
//             root1, root2 = root2, root1
//         self.parent[root1] = root2
//         self.rank[root2] += self.rank[root1]
//         self.count -= 1
//         return True

//     def find(self, key: T) -> T:
//         if key not in self.parent:
//             self.add(key)
//             return key

//         while self.parent.get(key, key) != key:
//             self.parent[key] = self.parent[self.parent[key]]
//             key = self.parent[key]
//         return key

//     def isConnected(self, key1: T, key2: T) -> bool:
//         return self.find(key1) == self.find(key2)

//     def getRoots(self) -> List[T]:
//         return list(set(self.find(key) for key in self.parent))

//     def getGroup(self) -> DefaultDict[T, List[T]]:
//         groups = defaultdict(list)
//         for key in self.parent:
//             root = self.find(key)
//             groups[root].append(key)
//         return groups

//     def add(self, key: T) -> bool:
//         if key in self.parent:
//             return False
//         self.parent[key] = key
//         self.rank[key] = 1
//         self.count += 1
//         return True

// 扩展域并查集判断二分图.
//  配合 `OfflineDynamicConnectivity` 可以支持动态图的二分图判断.
//  https://www.luogu.com.cn/problem/solution/P5787
func IsBipartite2() {}
