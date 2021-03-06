![树的重心](../../images/447bf9065511f5b136989797fc52f50f60c64eaf5b265a87c3fbf12f25c66860.png)
![树的中心](..%5C..%5C11_%E5%8A%A8%E6%80%81%E8%A7%84%E5%88%92%5Cacwingdp%E4%B8%93%E9%A1%B9%E7%BB%83%E4%B9%A0%5C%E6%A0%91%E5%BD%A2DP%5C1073.%20%E6%A0%91%E7%9A%84%E4%B8%AD%E5%BF%83.py)
![n 叉树的直径](..%5C..%5C11_%E5%8A%A8%E6%80%81%E8%A7%84%E5%88%92%5Cacwingdp%E4%B8%93%E9%A1%B9%E7%BB%83%E4%B9%A0%5C%E6%A0%91%E5%BD%A2DP%5C1072.%20%E6%A0%91%E7%9A%84%E6%9C%80%E9%95%BF%E8%B7%AF%E5%BE%84-n%E5%8F%89%E6%A0%91%E7%9A%84%E7%9B%B4%E5%BE%84.py)

```JS
function findCentre(n: number, edges: [next: number, weight: number][][]): number[] {
  const res: number[] = []
  // 最大连通块大小
  const maxSizeOfSubtree = Array<number>(n).fill(Infinity)
  // 树的大小
  const treeSize = Array<number>(n).fill(Infinity)

  function dfs(cur: number, parent: number) {
    treeSize[cur] = 1
    maxSizeOfSubtree[cur] = 0

    for (const [next, _] of edges[cur]) {
      if (next === parent) continue
      dfs(next, cur)
      // dfs更新阶段
      // 后序,更新cur:此时cur可以拿到各个next的信息，更新cur
      maxSizeOfSubtree[cur] = Math.max(maxSizeOfSubtree[cur], treeSize[next])
      treeSize[cur] += treeSize[next]
    }

    // dfs结算阶段:准备离开cur了，最后检查cur是否合法
    maxSizeOfSubtree[cur] = Math.max(maxSizeOfSubtree[cur], n - treeSize[cur])
    if (maxSizeOfSubtree[cur] <= n / 2) res.push(cur)
  }
  dfs(0, -1)

  console.log(treeSize, maxSizeOfSubtree)
  return res
}
```

后序 dfs 的两个地方：`要在脑海里想象递归树`

```JS
function dfs(cur: number, parent: number) {
    treeSize[cur] = 1
    maxSizeOfSubtree[cur] = 0

    for (const [next, _] of edges[cur]) {
      if (next === parent) continue
      dfs(next, cur)

      // 1....
      // 根据next的信息更新cur
    }

    // 2...
    // 准备撤离cur，结算cur
  }
```

1 那个地方是 cur 用他下面相邻的 next 结点信息更新自己
2 那个地方是 dfs 准备永远离开 cur 结点 做最后的收尾工作

## 树哈希(树的最小表示法/同构)

子树哈希为 `subtreeHash`
包含 root 的树哈希为 `${root.val}(${subtreeHash})`

```JS
function dfs(root: TrieNode): string {
  const subTree: string[] = []
  for (const child of root.children.values()) {
    subTree.push(dfs(child))
  }

  subTree.sort()
  root.subtreeHash = subTree.join('')
  hashCounter.set(root.subtreeHash, (hashCounter.get(root.subtreeHash) || 0) + 1)

  const res = `${root.val}(${root.subtreeHash})`
  return res
}
```

```Python
    def dfs(s: str, index: int) -> str:
        subtree = []
        depthDiff = 0
        for next in range(index + 1, len(s)):  # 找子树结点
            depthDiff += 1 if s[next] == '0' else -1
            if depthDiff == 1 and s[next] == '0':
                subtree.append(dfs(s, next))
            if depthDiff < 0:  # 回上面去了
                break

        # 子树最小表示
        subtree = deque(sorted(subtree))
        subtree.appendleft('(')
        subtree.append(')')

        # 返回整棵树
        subtree.appendleft('#')  # 当前元素
        return ''.join(subtree)
```
