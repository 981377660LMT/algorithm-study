一维线段树是二叉树
二维线段树是四叉树
三维线段树是八叉树
...

---

二维线段树（树状数组套动态开点线段树）
https://github.com/old-yan/CP-template/blob/a07b6fe0092e9ee890a0e35ada6ea1bb2c83ba05/DS/SegBIT.md#L1

二维线段树（动态开点线段树套动态开点线段树）
https://github.com/old-yan/CP-template/blob/a07b6fe0092e9ee890a0e35ada6ea1bb2c83ba05/DS/SegTree2d.md#L1
二维线段树可以实现二维树状数组的绝大部分功能，而且可以实现更强的功能，比如区域最值查询、区域按位与、区域按位或。`但是不能实现二维树状数组加强版的区域修改功能。这是因为二维树状数组本身也不能实现区域修改功能`，加强版的区域修改功能只不过是在加法运算符、差分性质下做到的特例。
**二维线段树的区间修改需要树套树实现**

树套树
https://ei1333.github.io/algorithm/segment-tree.html
https://ei1333.github.io/library/structure/segment-tree/segment-tree-2d-2.hpp
