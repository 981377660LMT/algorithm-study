https://www.zhihu.com/question/27840936
https://leer.moe/2018/12/17/data-structrue-treap/

什么是树堆（Treap）？

在 Treap 当中，维护平衡非常简单，只有一句话，就是通过维护小顶堆的形式来维持树的平衡
我们是通过维持堆的性质来保持平衡的，那么自然又会有一个新的问题。为什么维持堆的性质可以保证平衡呢?
答案很简单，因为我们在插入的时候，需要对每一个插入的 Node 随机附上一个 priority。堆就是用来维护这个 priority 的，保证树根一定拥有最小的 priority。正是由于这个 priority 是随机的，我们可以保证整棵树蜕化成线性的概率降到无穷低。
// 无法支持 islice 因为 Treap 结点随机旋转

java/C++ 的 set 都是不支持下标查询的
Treap 可以做到
