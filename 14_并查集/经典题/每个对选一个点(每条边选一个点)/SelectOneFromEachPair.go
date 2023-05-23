package main

func main() {

}

// 不能路径压缩

type SelectOneFromEachPair struct {
	Part int // 联通分量数
}

func NewSelectOneFromEachPair() *SelectOneFromEachPair {}

func (s *SelectOneFromEachPair) Union(u, v int) bool {}

func (s *SelectOneFromEachPair) Find(u int) int {}

func (s *SelectOneFromEachPair) Undo() int {}

// 联通分量为树的联通分量数.
func (s *SelectOneFromEachPair) CountTree() int {}

// 从每条边中恰好选一个点, 最多能选出多少个不同的点.
func (s *SelectOneFromEachPair) Solve() int {
	return s.Part - s.CountTree()
}
