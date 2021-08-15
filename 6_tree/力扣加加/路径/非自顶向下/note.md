非自顶向下一般用左右子树递归法
设计一个辅助函数 maxpath，调用自身求出以一个节点为根节点的左侧最长路径 left 和右侧最长路径 right，那么经过该节点的最长路径就是 left+right
接着只需要从根节点开始 dfs,不断比较更新全局变量即可

```JS
let res=0;
//以root为路径起始点的最长路径
function dfs(root) {
  if (!root) return 0;
  const left=maxPath(root.left);
  const right=maxPath(root.right);
  res = Math.max(res, left + right + root.val); //左右之和，更新全局变量
  return Math.max(left, right);   //返回左右路径较长者，对上面的贡献
}
return res
```

链接：https://leetcode-cn.com/problems/path-sum-iii/solution/yi-pian-wen-zhang-jie-jue-suo-you-er-cha-smch/
