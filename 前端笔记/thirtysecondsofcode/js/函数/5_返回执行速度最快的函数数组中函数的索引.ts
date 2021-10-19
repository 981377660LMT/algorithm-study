mostPerformant([
  () => {
    // Loops through the entire array before returning `false`
    ;[1, 2, 3, 4, 5, 6, 7, 8, 9, '10'].every(el => typeof el === 'number')
  },
  () => {
    // Only needs to reach index `1` before returning `false`
    ;[1, '2', 3, 4, 5, 6, 7, 8, 9, 10].every(el => typeof el === 'number')
  },
]) // 1

// 迭代次数越多，结果越可靠，但所需时间也越长。
function mostPerformant(funcs: ((...args: any[]) => any)[], iterations = 10000): number {
  const times = funcs.map(func => {
    const before = performance.now()
    for (let i = 0; i < iterations; i++) func()
    return performance.now() - before
  })

  return times.indexOf(Math.min(...times))
}
