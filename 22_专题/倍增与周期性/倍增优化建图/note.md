- 倍增拆分序列上的区间 `[start,end)`
  [DivideIntervalBinaryLift](DivideIntervalBinaryLift.go)
  [RangeUnionFindTreeOffline](RangeUnionFindTreeOffline.go)
- 倍增拆分倍增结构上的路径 `link(from,len)`
  [DividePathOnDoublingBinaryLift](DividePathOnDoublingBinaryLift.go)
- 倍增拆分树上路径 `path(from,to)`
  [DoublingLca32](DoublingLca32.go)
  [RangeToRangeGraphOnTree](RangeToRangeGraphOnTree.go)

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
     - 区间长度不为 2^i，**且操作满足幂等律(idempotent)**
       `[start,end) -> jump(start,k) + jump(end-2^k,k)`

---

一共拆分成[0,log]层，每层有 n 个元素.
`jumpId = level*n + index 表示第level层的第index个元素(0<=level<log+1,0<=index<n)`.

---

倍增结构与线段树的区别

1. 对于任意两段长度相等的区间，倍增的子树结构是相同的，而线段树的子树结构是不同的.
   倍增的子区间经过平移之后，仍然对应一个合法的子区间，而线段树的子区间平移之后不一定.
   基于倍增结构的这一特点，可以实现"区间并行操作".
2. 倍增结构在每一层代表的区间可以重叠，而线段树不能;
3. 倍增拆点专用于`倍增结构上的路径拆成 logn 个点`；
   线段树拆点专用于序列上的区间拆成 logn 个点，如果放到倍增结构(例如树)上，需要先利用树剖剖分成 logn 段序列，再对每段序列进行线段树拆点，时空复杂度为 `O(n*logn*logn)`

---

倍增优化建图
https://taodaling.github.io/blog/2020/03/18/binary-lifting/
