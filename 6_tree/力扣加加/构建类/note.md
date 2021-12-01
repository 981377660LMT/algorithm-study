普通二叉树的构建

1.  由 两种 dfs 构建
    有中序 在中序找根
    没中序(前+后) 找 postOrder[preOrder[1]] 即后序遍历中最后一个左子树中的元素
2.  由 bfs 构建
3.  其他

二叉搜索树的构建

构建类的模板

```JS
const root = new BinaryTree(...)
root.left = helper(...)
root.right = helper(...)
return root
```

前序+后序可以确定一棵树 层次单独就可以确定一棵树
