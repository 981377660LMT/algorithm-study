// 给定一个字符串 s 和一个整数 k 。你可以从 s 的前 k 个字母中选择一个，并把它加到字符串的末尾。
// 返回 在应用上述步骤的任意数量的移动后，字典上最小的字符串 。

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function orderlyQueue(s: string, k: number): string {
  if (k > 1) return s.split('').sort().join('')
  return findMinimunIsomorphic(s)
}

function findMinimunIsomorphic(str: string): string {
  if (str.length <= 1) return str

  const n = str.length
  let i = 0
  let j = 1
  let k = 0

  while (i < n && j < n && k < n) {
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    const diff = str.codePointAt((i + k) % n)! - str.codePointAt((j + k) % n)!

    if (diff === 0) {
      k++
      continue
    }

    if (diff > 0) i += k + 1
    else if (diff < 0) j += k + 1

    if (i === j) j++

    k = 0
  }

  const res = i > j ? j : i
  return str.slice(res) + str.slice(0, res)
}

export {}
