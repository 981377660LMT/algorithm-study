/* eslint-disable no-inner-declarations */
/* eslint-disable no-undef-init */
/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable implicit-arrow-linebreak */
/* eslint-disable no-await-in-loop */

type RecursiveGenerator<T> = Generator<T | RecursiveGenerator<T>, unknown, T>

/**
 * @param gen 递归生成器.
 * @param shouldBreak 是否中断.如果中断,则返回当前的结果.默认不中断.
 * @param useTimeSlicing 是否使用时间片.如果使用时间片,则每64ms让出cpu一次.默认不使用.
 * @returns 生成器的最终结果.
 */
async function bootStrapAsync<T>(
  gen: RecursiveGenerator<T>,
  shouldBreak?: () => boolean,
  useTimeSlicing?: boolean
): Promise<T> {
  const stack = [gen]
  let res: any = undefined
  let start = Date.now()

  while (stack.length) {
    if (shouldBreak && shouldBreak()) return res

    const last = stack[stack.length - 1]
    const next = last.next(res)
    if (next.done) {
      res = next.value
      stack.pop()
    } else {
      stack.push(next.value as RecursiveGenerator<T>)
    }

    if (useTimeSlicing) {
      const now = Date.now()
      if (now - start > 64) {
        await sleep()
        start = Date.now()
      }
    }
  }

  return res
}

function bootStrap<T>(gen: RecursiveGenerator<T>): T {
  const stack = [gen]
  let res = undefined as any
  while (stack.length) {
    const last = stack[stack.length - 1]
    const next = last.next(res)
    if (next.done) {
      res = next.value
      stack.pop()
    } else {
      stack.push(next.value as RecursiveGenerator<T>)
    }
  }

  return res
}

function sleep(time = 0): Promise<void> {
  return new Promise(resolve => {
    setTimeout(resolve, time)
  })
}

export { bootStrap, bootStrapAsync }

if (require.main === module) {
  // eslint-disable-next-line no-inner-declarations
  function facDfs(n: number): bigint {
    if (n <= 1) return BigInt(1)
    return BigInt(n) * facDfs(n - 1)
  }

  // eslint-disable-next-line no-inner-declarations
  function* facGen(n: number): RecursiveGenerator<bigint> {
    if (n <= 1) return BigInt(1)
    return BigInt(n) * (yield facGen(n - 1))
  }

  const adjList: number[][] = [[1, 2], [0, 2], [0, 1, 3], [2]]

  function dfs(cur: number, pre: number): void {
    console.log(cur)
    for (const next of adjList[cur]) {
      if (next !== pre) {
        dfs(next, cur)
      }
    }
  }

  console.log(facGen(10).next().value)

  const res = bootStrap(facGen(10))
  console.log(res)
}
