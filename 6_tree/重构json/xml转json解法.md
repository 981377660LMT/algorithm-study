与一维数组转树形数组解法类似,都是利用 dfs 的**副作用**

1. 准备好根节点的数组
2. dfs 时建立节点，将 新的 root push 到父节点，孩子节点同样 dfs
   **疑问：这里为什么要 root push 到父节点 而不是 将 root 子节点 push 到 root ?**
   因为将 root 子节点 push 到 root 需要 return
   **而 root push 到父节点不需要 return ,更加简洁**

```JS
一维数组转树形数组：
const dfs = (root: ArrayItem) => {
    if (adjMap.has(root.id)) {
      root.children = []
      for (const child of adjMap.get(root.id)!) {
        root.children.push(child)
        dfs(child)
      }
    }
  }

xml 转 json:
if (match.length) {
      for (const group of match) {
        ...
        const root = new TNode(name, [])
        parent.children!.push(root)
        dfs(childStr, root)
      }
    }
else {
    parent.children!.push(new TNode(str))
}
```
