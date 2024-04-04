/**
 * 对每个固定的左端点`left(0<=left<n)`，找到最大的右端点`maxRight`，
 * 使得滑动窗口内的元素满足`predicate(left,maxRight)`成立.
 * 如果不存在，`maxRight`为-1.
 */
function getMaxRight(
  n: number,
  operation: {
    append: (right: number) => void
    popLeft: (left: number) => void
    predicate: (left: number, right: number) => boolean
  }
): Int32Array {
  const maxRight = new Int32Array(n)
  let right = 0
  const visitedRight = new Uint8Array(n)
  for (let left = 0; left < n; left++) {
    if (right < left) right = left
    while (right < n) {
      if (!visitedRight[right]) {
        visitedRight[right] = 1
        operation.append(right)
      }
      if (operation.predicate(left, right)) {
        right++
      } else {
        break
      }
    }

    if (right === n) {
      maxRight.fill(n - 1, left)
      break
    }

    maxRight[left] = right - 1 >= left ? right - 1 : -1
    operation.popLeft(left)
  }

  return maxRight
}

export { getMaxRight }

if (require.main === module) {
  // 3. 无重复字符的最长子串
  // https://leetcode.cn/problems/longest-substring-without-repeating-characters/
  function lengthOfLongestSubstring(s: string): number {
    const n = s.length
    const counter = Array(128).fill(0)
    let dupCount = 0
    const maxRight = getMaxRight(n, {
      append: right => {
        const code = s.charCodeAt(right)
        if (counter[code] === 1) dupCount++
        counter[code]++
      },
      popLeft: left => {
        const code = s.charCodeAt(left)
        counter[code]--
        if (counter[code] === 1) dupCount--
      },
      predicate: () => dupCount === 0
    })
    let res = 0
    for (let left = 0; left < n; left++) {
      if (maxRight[left] === -1) break
      res = Math.max(res, maxRight[left] - left + 1)
    }
    return res
  }
}
