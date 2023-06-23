// https://runjs.co/s/WSzbfVXM2

type AsyncFunction<R> = (...args: unknown[]) => Promise<R>

/**
 * 计算`sum(range(start,end))%mod`.
 * 因为js是单线程的，所以这里的worker是无效的，实际上只有一个线程在运行。
 */
async function modSumAsync(start: number, end: number, mod: number, worker = 10): Promise<number> {
  if (start >= end) return 0
  const count = end - start
  const base = Math.floor(count / worker)
  const remain = count % worker

  const funcs: AsyncFunction<number>[] = Array(worker)
  for (let i = 0; i < worker; i++) funcs[i] = () => run(i)

  const res = await Promise.all(funcs.map(f => f()))
  return res.reduce((a, b) => (a + b) % mod, 0)

  async function run(workerId: number): Promise<number> {
    const more = Math.min(workerId, remain)
    const normal = Math.max(0, workerId - remain)
    const curStart = start + more * (base + 1) + normal * base
    const curEnd = curStart + base + +(more < remain)
    let sum = 0
    for (let i = curStart; i < curEnd; i++) sum = (sum + i) % mod
    return sum
  }
}

async function main() {
  // const time1 = performance.now()
  console.time('modSumAsync')
  await modSumAsync(0, 5e8, 1e11 + 7)
  console.timeEnd('modSumAsync')
  // console.log(performance.now() - time1) // 7359.80000000447
}

main()

export {}
