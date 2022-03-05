存信息

图删除和树删除的笔记
内部返回与外部更新是一样的 外部可以最后看到所有结果 存外面更好
两种返回的模式

1. dfs 本身返回，这种一般采用后序 dfs 例如返回[sum,count]:一般做法
2. 返回结果不在 dfs 中，dfs 返回 void，使用数组在外部全局记录
   例如 subtreeSum/subtreeCount
   其实这两种方式本质是一样的；`图遍历，知道节点 id 的话使用开数组更方便；树的话只需要在内部递归返回即可`
   `1519. 子树中标签相同的节点数.py`

3. 两个阶段

```JS
  function dfs(cur: number): void {
    for (const next of adjList[cur]) {
      dfs(next)

      // 1. 后序dfs，由这个分支更新
      subTreeSum[cur] += subTreeSum[next]
      subTreeCount[cur] += subTreeCount[next]
    }

    // 2. 统计完了 准备回溯
    if (subTreeSum[cur] === 0) subTreeCount[cur] = 0
  }
```
