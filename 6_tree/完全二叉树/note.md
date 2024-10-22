# 完美二叉树(Perfect Binary Tree)

完美二叉树 是指所有叶子节点都在同一层级的树，且每个父节点恰有两个子节点。

```python
def dfs(node: Optional[TreeNode]) -> Tuple[int, bool]:
    """(完美二叉树的高度,是否是完美二叉树)"""
    if not node:
        return 0, True
    leftHeight, isLeftPerfect = dfs(node.left)
    rightHeight, isRightPerfect = dfs(node.right)
    if isLeftPerfect and isRightPerfect and leftHeight == rightHeight:
        res.append((1 << (leftHeight + 1)) - 1)
        return leftHeight + 1, True
    return 0, False
```
