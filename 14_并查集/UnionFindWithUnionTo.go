// class UnionFindMapWithUnionTo(Generic[T]):
//     __slots__ = ("part", "_parent", "_rank")

//     def __init__(self, iterable: Optional[Iterable[T]] = None):
//         self.part = 0
//         self._parent = dict()
//         self._rank = dict()
//         for item in iterable or []:
//             self.add(item)

//     def union(self, key1: T, key2: T) -> bool:
//         """按秩合并."""
//         root1 = self.find(key1)
//         root2 = self.find(key2)
//         if root1 == root2:
//             return False
//         if self._rank[root1] > self._rank[root2]:
//             root1, root2 = root2, root1
//         self._parent[root1] = root2
//         self._rank[root2] += self._rank[root1]
//         self.part -= 1
//         return True

//     def unionTo(self, child: T, parent: T) -> bool:
//         """定向合并."""
//         root1 = self.find(child)
//         root2 = self.find(parent)
//         if root1 == root2:
//             return False
//         self._parent[root1] = root2
//         self._rank[root2] += self._rank[root1]
//         self.part -= 1
//         return True

//     def find(self, key: T) -> T:
//         if key not in self._parent:
//             self.add(key)
//             return key
//         while self._parent.get(key, key) != key:
//             self._parent[key] = self._parent[self._parent[key]]
//             key = self._parent[key]
//         return key

//     def isConnected(self, key1: T, key2: T) -> bool:
//         return self.find(key1) == self.find(key2)

//     def getRoots(self) -> List[T]:
//         return list(set(self.find(key) for key in self._parent))

//     def getGroups(self) -> DefaultDict[T, List[T]]:
//         groups = defaultdict(list)
//         for key in self._parent:
//             root = self.find(key)
//             groups[root].append(key)
//         return groups

//     def getSize(self, key: T) -> int:
//         return self._rank[self.find(key)]

//     def add(self, key: T) -> bool:
//         if key in self._parent:
//             return False
//         self._parent[key] = key
//         self._rank[key] = 1
//         self.part += 1
//         return True

//     def __repr__(self) -> str:
//         return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

//     def __contains__(self, key: T) -> bool:
//         return key in self._parent

// class UnionFindArrayWithUnionTo:
//     __slots__ = ("n", "part", "_parent", "_rank")

//     def __init__(self, n: int):
//         self.n = n
//         self.part = n
//         self._parent = list(range(n))
//         self._rank = [1] * n

//     def find(self, x: int) -> int:
//         while self._parent[x] != x:
//             self._parent[x] = self._parent[self._parent[x]]
//             x = self._parent[x]
//         return x

//     def union(self, x: int, y: int) -> bool:
//         """按秩合并."""
//         rootX = self.find(x)
//         rootY = self.find(y)
//         if rootX == rootY:
//             return False
//         if self._rank[rootX] > self._rank[rootY]:
//             rootX, rootY = rootY, rootX
//         self._parent[rootX] = rootY
//         self._rank[rootY] += self._rank[rootX]
//         self.part -= 1
//         return True

//     def unionTo(self, child: int, parent: int) -> bool:
//         """定向合并.将child的父节点设置为parent."""
//         rootX = self.find(child)
//         rootY = self.find(parent)
//         if rootX == rootY:
//             return False
//         self._parent[rootX] = rootY
//         self._rank[rootY] += self._rank[rootX]
//         self.part -= 1
//         return True

//     def unionWithCallback(self, x: int, y: int, f: Callable[[int, int], None]) -> bool:
//         """
//         f: 合并后的回调函数, 入参为 (big, small)
//         """
//         rootX = self.find(x)
//         rootY = self.find(y)
//         if rootX == rootY:
//             return False
//         if self._rank[rootX] > self._rank[rootY]:
//             rootX, rootY = rootY, rootX
//         self._parent[rootX] = rootY
//         self._rank[rootY] += self._rank[rootX]
//         self.part -= 1
//         f(rootY, rootX)
//         return True

//     def isConnected(self, x: int, y: int) -> bool:
//         return self.find(x) == self.find(y)

//     def getGroups(self) -> DefaultDict[int, List[int]]:
//         groups = defaultdict(list)
//         for key in range(self.n):
//             root = self.find(key)
//             groups[root].append(key)
//         return groups

//     def getRoots(self) -> List[int]:
//         return list(set(self.find(key) for key in self._parent))

//     def getSize(self, x: int) -> int:
//         return self._rank[self.find(x)]

// def __repr__(self) -> str:
//
//	return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())
package main

import "fmt"

func main() {
	uf := NewUnionFindArrayWithUnionTo(10)
	uf.UnionTo(1, 2)
	fmt.Println(uf.GetGroups())
}

type UnionFindArrayWithUnionTo struct {
	Part   int
	n      int
	parent []int32
	rank   []int32
}

func NewUnionFindArrayWithUnionTo(n int) *UnionFindArrayWithUnionTo {
	parent := make([]int32, n)
	rank := make([]int32, n)
	for i := 0; i < n; i++ {
		parent[i] = int32(i)
		rank[i] = 1
	}
	return &UnionFindArrayWithUnionTo{Part: n, n: n, parent: parent, rank: rank}
}

func (u *UnionFindArrayWithUnionTo) Find(x int) int {
	x32 := int32(x)
	for u.parent[x32] != x32 {
		u.parent[x32] = u.parent[u.parent[x32]]
		x32 = u.parent[x32]
	}
	return int(x32)
}

// 按秩合并.
func (u *UnionFindArrayWithUnionTo) Union(x, y int, f func(big, small int)) bool {
	rootX, rootY := u.Find(x), u.Find(y)
	if rootX == rootY {
		return false
	}
	if u.rank[rootX] > u.rank[rootY] {
		rootX, rootY = rootY, rootX
	}
	u.parent[rootX] = int32(rootY)
	u.rank[rootY] += u.rank[rootX]
	u.Part--
	if f != nil {
		f(rootY, rootX)
	}
	return true
}

// 定向合并.
func (u *UnionFindArrayWithUnionTo) UnionTo(child, parent int) bool {
	rootX, rootY := u.Find(child), u.Find(parent)
	if rootX == rootY {
		return false
	}
	u.parent[rootX] = int32(rootY)
	u.rank[rootY] += u.rank[rootX]
	u.Part--
	return true
}

func (u *UnionFindArrayWithUnionTo) GetSize(x int) int {
	return int(u.rank[u.Find(x)])
}

func (u *UnionFindArrayWithUnionTo) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for key := 0; key < u.n; key++ {
		root := u.Find(key)
		groups[root] = append(groups[root], key)
	}
	return groups
}

type UnionFindMapWithUnionTo struct {
	Part   int
	parent map[int32]int32
	rank   map[int32]int32
}

func NewUnionFindMapWithUnionTo() *UnionFindMapWithUnionTo {
	return &UnionFindMapWithUnionTo{Part: 0, parent: make(map[int32]int32), rank: make(map[int32]int32)}
}

func (u *UnionFindMapWithUnionTo) Find(key int) int {
	key32 := int32(key)
	if _, ok := u.parent[key32]; !ok {
		u.add(key32)
		return key
	}
	for u.parent[key32] != key32 {
		u.parent[key32] = u.parent[u.parent[key32]]
		key32 = u.parent[key32]
	}
	return int(key32)
}

// 按秩合并.
func (u *UnionFindMapWithUnionTo) Union(key1, key2 int, f func(big, small int)) bool {
	root1, root2 := int32(u.Find(key1)), int32(u.Find(key2))
	if root1 == root2 {
		return false
	}
	if u.rank[root1] > u.rank[root2] {
		root1, root2 = root2, root1
	}
	u.parent[root1] = root2
	u.rank[root2] += u.rank[root1]
	u.Part--
	if f != nil {
		f(int(root2), int(root1))
	}
	return true
}

// 定向合并.
func (u *UnionFindMapWithUnionTo) UnionTo(child, parent int) bool {
	root1, root2 := int32(u.Find(child)), int32(u.Find(parent))
	if root1 == root2 {
		return false
	}
	u.parent[root1] = root2
	u.rank[root2] += u.rank[root1]
	u.Part--
	return true
}

func (u *UnionFindMapWithUnionTo) GetSize(key int) int {
	return int(u.rank[int32(u.Find(key))])
}

func (u *UnionFindMapWithUnionTo) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for key := range u.parent {
		root := u.Find(int(key))
		groups[root] = append(groups[root], int(key))
	}
	return groups
}

func (u *UnionFindMapWithUnionTo) add(key32 int32) {
	u.parent[key32] = key32
	u.rank[key32] = 1
	u.Part++
}
