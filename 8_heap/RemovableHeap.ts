// 堆的删除思路有两种:
// !1. 一种是懒删除，即查询时再实际删除元素；
// https://github.com/harttle/contest.js/blob/master/src/heap.ts
// !2. 另一种是实时删除 index 处的元素，这种做法需要保证堆中的元素带有 index 属性。
// https://cs.opensource.google/go/go/+/refs/tags/go1.19.2:src/container/heap/heap.go
