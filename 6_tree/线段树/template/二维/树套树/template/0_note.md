树套树在日本竞赛圈叫 RangeTree(领域树)

参考 maspy/nyann/ei13333 的二维数树状数组和线段树的实现
https://nyaannyaan.github.io/library/data-structure-2d/abstract-range-tree.hpp

树状数组(离散/非离散)套数据结构
https://ei1333.github.io/library/structure/others/abstract-2d-binary-indexed-tree-compressed.hpp
https://ei1333.github.io/library/structure/others/abstract-binary-indexed-tree.hpp

线段树(非离散)套数据结构
https://ei1333.github.io/library/structure/segment-tree/segment-tree-2d-2.hpp

---

树套树维护维护多维度信息

https://oi-wiki.org/ds/seg-in-seg/

应用举例：

- 线段树套平衡树(SortedList) -> 区间前驱后继,第 k 大等
- 线段树套线段树 -> 一般的二维线段树
- 树状数组套权值线段树 -> O(nlogn)动态区间 kth，优于线段树套平衡树 O(nlognlogn)
- 分块套树状数组 -> 动态二维矩形区域查询；分块也可以看成一种树，因为 duck typing

面对多维度信息的题目时，如果题目没有要求强制在线，我们还可以考虑 CDQ 分治，或者 整体二分 等分治算法，来避免使用高级数据结构，减少代码实现难度。

---

golang 适合写，不容易 MLE
ts 可以作为设计指导

---

超冷门数据结构——二维线段树详解
https://www.luogu.com.cn/blog/Hoshino-kaede/chao-leng-men-shuo-ju-jie-gou-er-wei-xian-duan-shu-yang-xie

- 四叉树的错误复杂度->`n*1`的矩阵卡到最坏

---

**参考 divideInterval，树套树的本质就是利用线段树/树状数组 将区间分割成 logn 段，将区间修改与区间查询打到这些段上。**
