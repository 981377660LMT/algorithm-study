/**
 * @returns 返回 minLen <= 长度 <= maxLen 的回文数字
 */
function* genPalindromeByLength(
  minLen: number,
  maxLen: number,
  isReversed = false
): Generator<number, number, void> {
  if (isReversed) {
    for (let len = maxLen; len >= minLen; len--) yield* inner(len, isReversed)
  } else {
    for (let len = minLen; len <= maxLen; len++) yield* inner(len, isReversed)
  }

  return -1

  function* inner(length: number, isRevered = false): Generator<number, void, void> {
    const start = 10 ** ((length - 1) >> 1)
    const end = start * 10 - 1

    if (isRevered) {
      for (let half = end; half >= start; half--) {
        if (length & 1)
          yield Number(String(half).slice(0, -1) + String(half).split('').reverse().join(''))
        else yield Number(String(half) + String(half).split('').reverse().join(''))
      }
    } else {
      for (let half = start; half <= end; half++) {
        if (length & 1)
          yield Number(String(half).slice(0, -1) + String(half).split('').reverse().join(''))
        else yield Number(String(half) + String(half).split('').reverse().join(''))
      }
    }
  }
}

function kMirror(k: number, n: number): number {
  const res: number[] = []
  // 所有的回文数
  const iter = genPalindromeByLength(1, 1e9)

  while (res.length < n) {
    const palindrome = iter.next().value
    const kDigit = palindrome.toString(k)
    if (kDigit === kDigit.split('').reverse().join('')) {
      res.push(palindrome)
    }
  }

  return res.reduce((pre, cur) => pre + cur, 0)
}

export { genPalindromeByLength }

if (require.main === module) {
  console.log(kMirror(3, 7))
}
