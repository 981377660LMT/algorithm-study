1. 函数的入参全都是叫 root。而这个技巧是说，我们在写 dfs 函数的时候，要将函数中表示当前节点的形参也写成 root。
2. **自顶向下**就是在每个递归层级，首先访问节点来计算一些值，并在递归调用函数时将这些值传递到子节点，一般是通过参数传到子树中。
   **自底向上**是另一种常见的递归方法，首先对所有子节点递归地调用函数，然后根据返回值和根节点本身的值得到答案。
   大多数树的题使用**后序遍历比较简单，并且大多需要依赖左右子树的返回值**。比如 1448. 统计二叉树中好节点的数目
   不多的问题需要前序遍历，而**前序遍历通常要结合参数扩展技巧**。比如 1022. 从根到叶的二进制数之和
   如果你能使**用参数和节点本身的值来决定什么应该是传递给它子节点的参数，那就用前序遍历**。
   如果对于树中的任意一个节点，如果你**知道它子节点的答案，你能计算出当前节点的答案，那就用后序遍历**。
   如果遇到**二叉搜索树则考虑中序遍历**

3. 删除节点的题目:
   - 450. 删除二叉搜索树中的节点.ts
   - 669. 修剪二叉搜索树.ts
   - 814. 二叉树剪枝.ts
   - 1325 删除给定值的叶子节点.ts
     父节点可能被删除
     一般是后序遍历+递归

```JS
  root.left && (root.left = pruneTree(root.left))
  root.right && (root.right = pruneTree(root.right))
  // 主要的逻辑
```

4. 边界:空节点/叶子节点
5. 扩展 dfs 参数

   - 携带父亲或者爷爷的信息

   ```Python
   def dfs(root, parent):
      if not root: return
      dfs(root.left, root)
      dfs(root.right, root)
   ```

   - 携带路径信息，可以是路径和或者具体的路径数组等

   ```Python
   def dfs(root, path_sum):
    if not root:
        # 这里可以拿到根到叶子的路径和
        return path_sum
    dfs(root.left, path_sum + root.val)
    dfs(root.right, path_sum + root.val)
   ```

   ```Python
   def dfs(root, path):
    if not root:
        # 这里可以拿到根到叶子的路径
        return path
    path.append(root.val)
    dfs(root.left, path)
    dfs(root.right, path)
    # 撤销
    path.pop()
   ```

   - 二叉搜索树的搜索题大多数都需要扩展参考，甚至怎么扩展都是固定的
     二叉搜索树的搜索总是将最大值和最小值通过参数传递到左右子树，类似 dfs(root, lower, upper)，然后在递归过程更新最大和最小值即可。

6. dfs 函数的返回值

- 1530. 好叶子节点对的数量.ts dfs 返回距离数组
- 894. 所有可能的满二叉树
