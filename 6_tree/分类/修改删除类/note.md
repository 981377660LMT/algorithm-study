题目要求的修改

- 修改节点的值或者指向
- 增加，删除节点

算法需要，自己修改

1. 删除一般是 dfs 后序

```JS
    root.left = dfs(root.left, 删除条件/标志位)
    root.right = dfs(root.right, 删除条件/标志位)
    if(...) return null
    return root
```

2. 使用虚拟节点(类似链表删除)
   `814. 二叉树剪枝.py`
   `1080. 根到叶路径上的不足节点.ts`
   `1325. 删除给定值的叶子节点.ts`

```JS
function pruneTree(root: BinaryTree | null): BinaryTree | null {
  const dummy = new BinaryTree(0, root)
  dfs(dummy)
  return dummy.left

  function dfs(root: BinaryTree | null): boolean {
    if (!root) return false
    const left = dfs(root.left)
    const right = dfs(root.right)
    if (!root.left) {
        root.left=null
    }
    ...
  }
}

```
