题目要求的修改

- 修改节点的值或者指向
- 增加，删除节点

算法需要，自己修改

删除一般是 dfs 前序 递归删除

```JS
    root.left = dfs(root.left, 删除条件/标志位)
    root.right = dfs(root.right, 删除条件/标志位)
    return root
```
