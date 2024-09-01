## 树的欧拉序列(euler tour)

对一棵树进行 DFS，无论是第一次访问还是回溯，每次到达一个结点时都将编号记录下来，可以得到一个长度为 `2n-1 (2*(n-1)+1)` 的序列，这个序列被称作这棵树的欧拉序列。

```Python
def dfs(cur: int, pre: int, eulerTour: List[int]) -> None:
    eulerTour.append(cur)  # 访问
    for next in adjList[cur]:
        if next != pre:
            dfs(next, cur, eulerTour)
            eulerTour.append(cur)  # 回溯
```

## 树的欧拉序列的应用 (应用少，不如重链剖分处理路径通用)

https://maspypy.com/euler-tour-%e3%81%ae%e3%81%8a%e5%8b%89%e5%bc%b7

理解：
每条路径可以拆成到 lca 的`两条链`，而`每条链`在欧拉序列中都是`连续的`，所以可以用`区间`来表示一条路径。
**配合欧拉序列的前缀和**，可以查询出链上的边权/点权和。

- 具有逆元的 monoid (群)可以由欧拉序列 O(logn)查询 **(因为在回溯时需要在前缀和中加上逆元)**
- 不具有逆元的 monoid (幺半群)可以由重链剖分 `O(logn*logn)`查询

## `访问指定k个节点的最短路径`

专题
