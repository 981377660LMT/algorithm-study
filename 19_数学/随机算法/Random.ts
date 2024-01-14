/* eslint-disable no-inner-declarations */

type Sequence<T> = { length: number; [index: number]: T }

class Random {
  static randInt(min: number, max: number): number {
    return min + Math.floor((max - min + 1) * Math.random())
  }

  static randFloat(min: number, max: number): number {
    return min + (max - min) * Math.random()
  }

  static randRange(start: number, stop: number, step = 1): number {
    const width = stop - start
    if (step === 1) return start + Random.randInt(0, width - 1)
    let n = 0
    if (step > 0) {
      n = Math.floor((width + step - 1) / step)
    } else {
      n = Math.floor((width + step + 1) / step)
    }
    return start + step * Random.randInt(0, n - 1)
  }

  /** Fisher-Yates. */
  static shuffle<T>(arr: Sequence<T>): void {
    for (let i = 0; i < arr.length; i++) {
      const rand = Random.randInt(0, i)
      Random._swap(arr, i, rand)
    }
  }

  static sample<T>(arr: Sequence<T>, k = 1): T[] {
    const copy = Array(arr.length)
    for (let i = 0; i < arr.length; i++) copy[i] = arr[i]
    Random.shuffle(copy)
    return copy.slice(0, k)
  }

  private static _swap<T>(arr: Sequence<T>, i: number, j: number): void {
    const tmp = arr[i]
    arr[i] = arr[j]
    arr[j] = tmp
  }
}

export { Random }

if (require.main === module) {
  const counter = new Map<number, number>()
  for (let _ = 0; _ < 100; _++) {
    const rand = Random.randRange(1, 10, 2)
    counter.set(rand, (counter.get(rand) ?? 0) + 1)
  }
  console.log(counter)

  // 给你一个下标从 0 开始的字符串 s 和一个整数 k。
  // 你需要执行以下分割操作，直到字符串 s 变为 空：
  // 选择 s 的最长前缀，该前缀最多包含 k 个 不同 字符。
  // 删除 这个前缀，并将分割数量加一。如果有剩余字符，它们在 s 中保持原来的顺序。
  // 执行操作之 前 ，你可以将 s 中 至多一处 下标的对应字符更改为另一个小写英文字母。
  // 在最优选择情形下改变至多一处下标对应字符后，用整数表示并返回操作结束时得到的最大分割数量。

  // 100154. 执行操作后的最大分割数量
  // https://leetcode.cn/problems/maximize-the-number-of-partitions-after-operations/description/
  // 随机抽取500个点，暴力枚举每个点的改变，然后计算最大分割数量
  // O(500*1e4*26)
  function maxPartitionsAfterOperations(s: string, k: number): number {
    const n = s.length
    const ords = new Uint8Array(n)
    for (let i = 0; i < s.length; i++) ords[i] = s[i].charCodeAt(0) - 97

    let res = 0
    const range = Array.from({ length: n }, (_, i) => i)
    const points = Random.sample(range, Math.min(1000, n))

    for (let i = 0; i < points.length; i++) {
      const changePoint = points[i]
      const pre = ords[changePoint]
      for (let v = 0; v < 26; v++) {
        ords[changePoint] = v
        let curRes = 0
        let ptr = 0
        while (ptr < n) {
          let visited = 1 << ords[ptr]
          let visitedCount = 1
          ptr++
          while (ptr < n && visitedCount + (1 ^ ((visited >>> ords[ptr]) & 1)) <= k) {
            visitedCount += ((visited >>> ords[ptr]) & 1) ^ 1
            visited |= 1 << ords[ptr]
            ptr++
          }
          curRes++
        }
        res = Math.max(res, curRes)
      }
      ords[changePoint] = pre
    }
    return res
  }
}
