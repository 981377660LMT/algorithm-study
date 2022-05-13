回溯是 DFS 中的一种技巧。回溯法采用 试错 的思想，它尝试分步的去解决一个问题。在分步解决问题的过程中，当它通过尝试发现现有的分步答案不能得到有效的正确的解答的时候，它将**取消上一步甚至是上几步的计算，再通过其它的可能的分步解答再次尝试寻找问题的答案。**

```JS
const visited = {}
function dfs(i) {
    if (满足特定条件）{
        // 返回结果 or 退出搜索空间
    }

    visited[i] = true // 将当前状态标为已搜索
    dosomething(i) // 对i做一些操作
    for (根据i能到达的下个状态j) {
        if (!visited[j]) { // 如果状态j没有被搜索过
            dfs(j)
        }
    }
    undo(i) // 恢复i
}
```

回溯题目的另外一个考点是剪枝， 通过恰当地剪枝，可以有效减少时间
避免根本不可能是答案的递归

注意：
在递归到底部往上冒泡的时候进行撤销状态。
**如果你每次递归的过程都拷贝了一份数据，那么就不需要撤销状态，相对地空间复杂度会有所增加**。
一些回溯的题目，我们仍然也可以采用**笛卡尔积**的方式，**将结果保存在返回值而不是路径引用中**，这样就避免了回溯状态，并且由于结果在返回值中，因此可以使用记忆化递归， 进而优化为动态规划形式。

字符串拆分：一般是 next=remain.slice(0, i + 1) new=remain.slice(i + 1)

```JS

  const bt = (path: string[], remain: string) => {
    if (remain.length === 0) {
      res.push(path.join(' '))
      return path.pop()
    }

    for (let i = 0; i < remain.length; i++) {
      const next = remain.slice(0, i + 1)
      if (store.has(next)) {
        path.push(next)
        bt(path, remain.slice(i + 1))
      }
    }

    path.pop()
  }
  bt([], s)
```
