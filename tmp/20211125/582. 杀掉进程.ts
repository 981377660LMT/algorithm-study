// 每一个进程只有 一个父进程 ，但是可能会有 一个或者多个子进程
// 只有一个进程的 ppid[i] = 0 ，意味着这个进程 没有父进程 。
// 给你一个整数 kill 表示要杀掉​​进程的 ID ，返回杀掉该进程后的所有进程 ID 的列表。可以按 任意顺序 返回答案。
// 1 <= n <= 5 * 104
// pid 中的所有值 互不相同

function killProcess(pid: number[], ppid: number[], kill: number): number[] {
  const n = pid.length
  const adjMap = new Map<number, number[]>()

  for (let i = 0; i < n; i++) {
    const cur = ppid[i]
    const next = pid[i]
    !adjMap.has(cur) && adjMap.set(cur, [])
    adjMap.get(cur)!.push(next)
  }

  const res: number[] = []
  dfs(kill)
  return res

  function dfs(root: number) {
    res.push(root)
    for (const next of adjMap.get(root) || []) {
      dfs(next)
    }
  }
}

console.log(killProcess([1, 3, 10, 5], [3, 0, 5, 3], 5))
// 输出：[5,10]
// 解释：涂为红色的进程是应该被杀掉的进程。
