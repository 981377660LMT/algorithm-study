// 分治的迭代写法

function mergeAll<T>(items: ArrayLike<T>, e: () => T, op: (e1: T, e2: T) => T): T {
  if (!items.length) return e()

  const copy = Array(items.length)
  for (let i = 0; i < items.length; i++) copy[i] = items[i]
  let n = copy.length
  while (n > 1) {
    const mid = (n + 1) >>> 1
    for (let i = 0; i < mid; i++) {
      if (((i << 1) | 1) ^ n) {
        copy[i] = op(copy[i << 1], copy[(i << 1) | 1])
      } else {
        copy[i] = copy[i << 1]
      }
    }
    n = mid
  }

  return copy[0]
}

async function mergeAllAsync<T>(
  items: Array<Promise<T>>,
  e: () => Promise<T>,
  op: (e1: Promise<T>, e2: Promise<T>) => Promise<T>
): Promise<T> {
  if (!items.length) return e()

  const copy = items.slice()
  let n = copy.length
  while (n > 1) {
    const mid = (n + 1) >> 1
    for (let i = 0; i < mid; i++) {
      // "!==n"
      if (((i << 1) | 1) ^ n) {
        copy[i] = op(copy[i << 1], copy[(i << 1) | 1])
      } else {
        copy[i] = copy[i << 1]
      }
    }
    n = mid
  }

  return copy[0]
}

if (require.main === module) {
  const sum = mergeAll(
    [1, 2, 3, 4, 5],
    () => 0,
    (a, b) => a + b
  )
  console.log(sum)
}

export { mergeAll, mergeAllAsync }
