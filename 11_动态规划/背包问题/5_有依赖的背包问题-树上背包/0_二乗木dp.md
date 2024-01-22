https://blog.hamayanhamayan.com/entry/2017/06/19/161741

- 二乗の木 DP という、頂点集合の DP をマージする時に部分木の要素数の個数分だけ使ってマージするようにすると O(N^3)が O(N^2)に落ちるテクがある

复杂度为平方的一系列树 dp，合并子树时**利用子树的大小**，将复杂度从 O(N^3) 降到 O(N^2)。
当然卷积合并也可以 多一个 log

---

树上背包优化：`二乘木dp`

https://zhuanlan.zhihu.com/p/316010761

https://yukicoder.me/problems/no/196
https://www.luogu.com.cn/problem/solution/CF815C
