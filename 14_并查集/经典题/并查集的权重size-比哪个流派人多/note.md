- 924. 尽量减少恶意软件的传播(`取消点感染病毒`)

https://leetcode.cn/problems/minimize-malware-spread/description/

求只包含一个被感染节点的最大连通块

改成移除 k 个节点呢？
因为每个连通块的感染数要么全部移除，要么全部保留。
移除 k 个节点可以在算出联通块的大小和每个组内被感染的个数之后采用 `01背包`的作法。
将连通块大小看成价值，每个块里面的感染个数是重量，而 k 就是背包的总容量。

- 928. 尽量减少恶意软件的传播 II (`从图中删除点`)

https://leetcode.cn/problems/minimize-malware-spread-ii/description/
