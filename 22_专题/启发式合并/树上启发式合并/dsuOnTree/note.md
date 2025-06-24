解决这样一类树上统计问题：

- 无修改操作，询问允许离线
- 对子树信息进行统计（链上的信息在某些条件下也可以统计）

**比一般的启发式合并快一些(使用空间少，更快)，作用是一样的。**
https://hitonanode.github.io/cplib-cpp/tree/guni.hpp
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

---

https://blog.csdn.net/qq_43472263/article/details/104150940
然后我们可以考虑一个优化，遍历到最后一个子树时是不用清空的，因为它不会产生对其他节点影响了，根据贪心的思想我们当然要把节点数最多的子树(即重儿子形成的子树)放在最后，之后我们就有了一个看似比较快的算法，**先遍历所有的轻儿子节点形成的子树，统计答案但是不保留数据，然后遍历重儿子，统计答案并且保留数据，最后再遍历轻儿子以及父节点，合并重儿子统计过的答案。**

---

启发式合并：size 小的合并到 size 大的
dsu on tree ：根据轻重链添加、删除子树贡献 https://www.luogu.com.cn/article/6028aw8w
