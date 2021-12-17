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
