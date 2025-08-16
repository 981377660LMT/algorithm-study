/**
 * 在 [start, end) 区间内，寻找第一个和最后一个满足 x % mod == remainder 的数。
 * 如果不存在，返回 [undefined, undefined]。
 */
function firstLastMod(
  start: number,
  end: number,
  mod: number,
  remainder: number
): [number, number] | [undefined, undefined] {
  if (start >= end) return [undefined, undefined]
  if (remainder < 0 || remainder >= mod) return [undefined, undefined]
  const r = start % mod
  const delta = (remainder - r + mod) % mod
  const first = start + delta
  if (first >= end) return [undefined, undefined]
  // 从 first 开始，每次加 mod，直到不超过 end-1
  const last = first + Math.floor((end - 1 - first) / mod) * mod
  return [first, last]
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  // 测试用例
  console.assert(JSON.stringify(firstLastMod(1, 10, 3, 2)) === JSON.stringify([2, 8]))
  console.assert(JSON.stringify(firstLastMod(1, 10, 3, 1)) === JSON.stringify([1, 7]))
  console.assert(JSON.stringify(firstLastMod(1, 11, 3, 1)) === JSON.stringify([1, 10]))
  console.assert(JSON.stringify(firstLastMod(1, 10, 3, 0)) === JSON.stringify([3, 9]))
  console.assert(
    JSON.stringify(firstLastMod(0, 0, 3, 0)) === JSON.stringify([undefined, undefined])
  )
  console.assert(
    JSON.stringify(firstLastMod(5, 5, 3, 2)) === JSON.stringify([undefined, undefined])
  )
  console.assert(JSON.stringify(firstLastMod(10, 20, 1, 0)) === JSON.stringify([10, 19]))
  console.assert(JSON.stringify(firstLastMod(7, 8, 2, 1)) === JSON.stringify([7, 7]))
  console.assert(
    JSON.stringify(firstLastMod(7, 8, 2, 0)) === JSON.stringify([undefined, undefined])
  )
  console.assert(
    JSON.stringify(firstLastMod(1, 10, 3, 3)) === JSON.stringify([undefined, undefined])
  )
  console.assert(
    JSON.stringify(firstLastMod(1, 10, 3, -1)) === JSON.stringify([undefined, undefined])
  )
  console.log('All tests passed.')
}

export { firstLastMod }
