type F = () => Promise<unknown>

// 限制promise并发数
function promisePool(functions: F[], n: number): Promise<unknown> {
  const res: unknown[] = []
  const iter = functions.entries()
  const workers = Array(n).fill(iter).map(work)
  return Promise.all(workers).then(() => res)

  async function work(entries: IterableIterator<[number, F]>) {
    for (const [index, task] of entries) {
      // eslint-disable-next-line no-await-in-loop
      const cur = await task()
      res[index] = cur
    }
  }
}

/**
 * const sleep = (t) => new Promise(res => setTimeout(res, t));
 * promisePool([() => sleep(500), () => sleep(400)], 1)
 *   .then(console.log) // After 900ms
 */

export {}
