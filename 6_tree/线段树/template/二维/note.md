一维线段树是二叉树
二维线段树是四叉树
三维线段树是八叉树
...

1. 单点修改，区间查询 => 一般的线段树

   二维线段树可以实现二维树状数组的绝大部分功能，而且可以实现更强的功能，比如区域最值查询、区域按位与、区域按位或。`但是不能实现二维树状数组加强版的区域修改功能。这是因为二维树状数组本身也不能实现区域修改功能`，加强版的区域修改功能只不过是在加法运算符、差分性质下做到的特例。
   **二维线段树的区间修改需要树套树实现**

2. 区间修改，单点查询 => 树套树
   树套树`无法做到一般的区间查询`
3. 区间修改，区间查询 => 递归版线段树(四叉树)
