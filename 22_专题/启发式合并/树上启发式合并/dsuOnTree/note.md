解决这样一类树上统计问题：

- 无修改操作，询问允许离线
- 对子树信息进行统计（链上的信息在某些条件下也可以统计）

比一般的启发式合并快一些，作用是一样的。
https://pzy.blog.luogu.org/dsu-on-tree-xue-xi-bi-ji
https://oi-wiki.org/graph/dsu-on-tree/
https://nyaannyaan.github.io/library/tree/dsu-on-tree.hpp

- 遍历所有轻儿子，递归结束时消除它们的贡献
- 遍历所有重儿子，保留它的贡献
- 再计算当前子树中所有轻子树的贡献
- 更新答案
- 如果当前点是轻儿子，消除当前子树的贡献

```py
# c ... current node
# p ... parent node
# keep ... condition variable of reserving data

def dsu(c, p, keep):
    # light edge -> run dfs and clear data
    for d in 'light edge of c':
        dsu(d, c, false)

    # heavy edge -> run dfs and reserve data
    dsu('heavy edge of c', c, true)

    # light edge -> reserve data
    for d in 'light edge of c':
        for n in 'subtree of d':
            add(n)

    # current node -> reserve data
    add(c)

    # answer queries related to subtree of current node
    query(c)

    # if keep is false, clear all data
    if keep = false:
          reset()
    return
```
