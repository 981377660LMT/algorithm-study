- 倍增拆分序列上的区间 `[start,end)`
  [DivideIntervalBinaryLift](DivideIntervalBinaryLift.go)
  [RangeUnionFindTreeOffline](RangeUnionFindTreeOffline.go)
- 倍增拆分倍增结构上的路径 `link(from,len)`
  [DividePathOnDoublingBinaryLift](DividePathOnDoublingBinaryLift.go)
- 倍增拆分树上路径 `path(from,to)`
  [DividePathOnTreeBinaryLift](DividePathOnTreeBinaryLift.go)

---

拆区间方法：

1. push up

   - 倍增结构
     `jump(from,i+1) <- jump(from,i) + jump(jump(from,i),i)`
   - 序列(st 表)
     `jump(from,i+1) <- jump(from,i) + jump(from+2^i,i)`

2. push down
   - 倍增结构
     `jump(start,i+1) -> jump(start,i) + jump(jump(start,i),i)`
   - 序列(st 表)
     - 区间长度为 2^i
       `jump(start,k+1) -> jump(start,k) + jump(start+2^k,k)`
     - 区间长度不为 2^i
       `[start,end) -> jump(start,k) + jump(end-2^k,k)`

---

一共拆分成[0,log]层，每层有 n 个元素.
`jumpId = level*n + index 表示第level层的第index个元素(0<=level<log+1,0<=index<n)`.
