ps:

1. 将后序遍历的结果进行反转，就是拓扑排序的结果。
   后序遍历的这一特点很重要，之所以拓扑排序的基础是后序遍历，
   **是因为一个任务必须在等到所有的依赖任务都完成之后才能开始开始执行。**
2. 拓扑排序只能针对有向无环图，进行拓扑排序之前要进行环检测

拓扑排序模板

```JS

  const inDegrees = Array<number>(n).fill(0)
  const adjList = Array.from<unknown, number[]>({ length: n }, () => [])
  for (const [cur, pre] of prerequisites) {
    inDegrees[cur]++
    adjList[pre].push(cur)
  }

  const queue: number[] = []
  inDegrees.forEach((v, i) => !v && queue.push(i))

  let count = 0
  while (queue.length) {
    const cur = queue.shift()!
    count++
    for (const next of adjList[cur]) {
      inDegrees[next]--
      if (inDegrees[next] === 0) queue.push(next)
    }
  }

  return count === n
```

3. 拓扑排序最短路径
   `2050. 并行课程 III-拓扑排序最短路径.py`
