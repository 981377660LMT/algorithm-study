// 堆的删除思路有两种:
// !1. 一种是懒删除，即查询时再实际删除元素；
// https://github.com/harttle/contest.js/blob/master/src/heap.ts
// !2. 另一种是实时删除 index 处的元素
// 调用 heappush 会返回一个 *viPair 指针，记作 p
// 将 p 存于他处（如 slice 或 map），可直接在外部修改 p.v 后调用 fix(p.index)，从而做到修改堆中指定元素
// !调用 remove(p.index) 可以从堆中删除 p.v
// https://github.dev/EndlessCheng/codeforces-go/blob/6be496d4d93d667e718f7f3db5519139a5f17ddf/copypasta/heap.go#L94
// https://cs.opensource.google/go/go/+/refs/tags/go1.19.2:src/container/heap/heap.go
