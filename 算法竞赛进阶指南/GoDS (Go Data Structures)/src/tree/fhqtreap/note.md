比较常用的平衡树一般就是 fhqtreap 和 splay

https://www.luogu.com.cn/blog/Chanis/fhq-treap

https://zhuanlan.zhihu.com/p/400448696
https://www.luogu.com.cn/problem/solution/P5055

1. treap 是一棵笛卡尔树，可以 O(n)时间建树。
2. 当树的结构发生改变的时候，当我们进行分裂或合并操作时需要改变某一个点的左右儿子信息时之前，应该下放标记(pushDown)
3. 当一个节点的左右儿子发生的变化的时候，需要上传维护值(pushUp)，在 split 和 merge 的时候就是一个例子
4. 持久化：为了**保证每次操作对历史版本不产生影响**，在 PushDown()和 Split()的操作时复制节点。因为每次 Merge()前都会 Split()，所以 Merge()时就不需要复制节点了。
5. 如果有区间操作，必须用排名分裂
6. Split 和 MergeMerge 操作都只会改变两颗 TreapTreap 交接处的节点

7. 区间操作时，只需分成三段，然后`中间那段当成线段树来操作就可以了`，注意区间更新的顺序
