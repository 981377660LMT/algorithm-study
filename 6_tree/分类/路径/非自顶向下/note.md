需要在`每个结点`处更新的树形 dp

1. 树形 dp，在`每个节点处更新全局的 res`,后序 dfs 返回子树内的信息
2. 记忆化 dfs ，答案为从`每个结点出发时`的 res 最大值

```JS
let res=0

// 以root为路径起始点的最长路径
function dfs(root) {
  if (!root) return 0;
  const left=maxPath(root.left)
  const right=maxPath(root.right)
  res = Math.max(res, left + right + root.val) // 左右之和，更新全局变量
  return Math.max(left, right)   // 返回左右路径较长者，对上面的贡献
}

return res
```
