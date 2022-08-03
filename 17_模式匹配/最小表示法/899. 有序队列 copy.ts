// 给定一个字符串 s 和一个整数 k 。你可以从 s 的前 k 个字母中选择一个，并把它加到字符串的末尾。
// 返回 在应用上述步骤的任意数量的移动后，字典上最小的字符串 。

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function orderlyQueue(s: string, k: number): string {
  if (k > 1) return s.split('').sort().join('')
  return findIsomorphic(s)
}

function findIsomorphic(str: string, isMin = true): string {
  if (str.length <= 1) return str

  const n = str.length
  let i1 = 0
  let i2 = 1
  let same = 0

  while (i1 < n && i2 < n && same < n) {
    const diff = compare(str[(i1 + same) % n], str[(i2 + same) % n])

    if (diff === 0) {
      same++
      continue
    }

    if (diff > 0) i1 += same + 1
    else if (diff < 0) i2 += same + 1

    if (i1 === i2) i2++

    same = 0
  }

  const res = Math.min(i1, i2)

  return `${str.slice(res)}${str.slice(0, res)}`

  function compare(a: string, b: string): number {
    if (a === b) return 0
    if (isMin) return a > b ? 1 : -1
    return a < b ? -1 : 1
  }
}

export {}
