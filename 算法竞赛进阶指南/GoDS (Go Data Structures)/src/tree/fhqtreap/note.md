https://www.luogu.com.cn/problem/solution/P5055 (深度好文)

比较常用的平衡树一般就是 fhqtreap 和 splay

https://www.luogu.com.cn/blog/Chanis/fhq-treap
https://zhuanlan.zhihu.com/p/400448696

1. treap 是一棵笛卡尔树，可以 O(n)时间建树。
2. 当树的结构发生改变的时候，当我们进行分裂或合并操作时需要改变某一个点的左右儿子信息时之前，应该下放标记(pushDown)
3. 当一个节点的左右儿子发生的变化的时候，需要上传维护值(pushUp)，在 split 和 merge 的时候就是一个例子
4. 持久化：为了**保证每次操作对历史版本不产生影响**，在 PushDown()和 Split()的操作时复制节点。因为每次 Merge()前都会 Split()，所以 Merge()时就不需要复制节点了。
5. Split 和 MergeMerge 操作都只会改变两颗 TreapTreap 交接处的节点
6. 区间操作时，只需分成三段，然后`中间那段当成线段树来操作就可以了`，注意区间更新的顺序
7. 持久化操作：Split 过程中可以对点进行复制，并且每次修改的必然只有一个子树上的点。而且 Split 和 Merge 总是成对出现，我们就只用复制一次。

8. 因为 Treap 是 BST，所以中序遍历按照(排名或者值)从小到大遍历，也就得到了整个 Treap 本身(有序序列)。
9. 分裂 split 有两种：排名(index)分裂、值(value)分裂；如果有区间操作，必须用排名分裂
   按照值分裂 => 名次树、SortedList 系列
   https://nyaannyaan.github.io/library/rbst/treap.hpp

   按照排名分裂 => 数组 系列
   https://nyaannyaan.github.io/library/rbst/lazy-reversible-rbst.hpp
   https://www.acwing.com/solution/content/51762/
   https://www.acwing.com/activity/content/code/content/2479443/

---

SortedList 做不到 O(1)寻找前驱后继, TreeSet 可以

迭代器的 api 设计
