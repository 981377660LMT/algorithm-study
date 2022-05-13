模板
双递归的基本套路就是一个主递归函数和一个内部递归函数。主递归函数负责计算以某一个节点开始的 xxxx，内部递归函数负责计算 xxxx，这样就实现了以所有节点开始的 xxxx。

```Python
def dfs_main(root):
  ## 以root为根的子问题
  def dfs_inner(root):
    # 这里写你的逻辑，就是前序遍历
    dfs_inner(root.left)
    dfs_inner(root.right)
    # 或者在这里写你的逻辑，那就是后序遍历
  ## 所有节点为根
  return dfs_inner(root) + dfs_main(root.left) + dfs_main(root.right)
```

是否存在一条子路径之和等于目标和
子路径 不需要从根节点开始，也不需要在叶子节点结束，但是路径方向必须是向下的（只能从父节点到子节点）。
若路径方向不定 则要用后序遍历 dfs 返回多个值`124.二叉树中的最大路径和`
