/**
 * 对每个固定的右端点`right(0<=right<n)`，找到最小的左端点`minLeft`，
 * 使得滑动窗口内的元素满足`predicate(minLeft,right)`成立.
 * 如果不存在，`minLeft`为`n`.
 */
function getMinLeft(
  n: number,
  operation: {
    append: (right: number) => void
    popLeft: (left: number) => void
    predicate: (left: number, right: number) => boolean
  }
): Int32Array {
  const minLeft = new Int32Array(n)
  let left = 0
  for (let right = 0; right < n; right++) {
    operation.append(right)
    while (left <= right && !operation.predicate(left, right)) {
      operation.popLeft(left)
      left++
    }
    minLeft[right] = left > right ? n : left
  }
  return minLeft
}

export { getMinLeft }

if (require.main === module) {
  const nums = [1, 2, 3, 4, 5]
  const n = nums.length
  let curSum = 0
  const minLeft = getMinLeft(n, {
    append: right => {
      curSum += nums[right]
    },
    popLeft: left => {
      curSum -= nums[left]
    },
    predicate: () => curSum <= 4
  })
  console.log(minLeft)

  // 3. 无重复字符的最长子串
  // https://leetcode.cn/problems/longest-substring-without-repeating-characters/
  function lengthOfLongestSubstring(s: string): number {
    const n = s.length
    const counter = Array(128).fill(0)
    let dupCount = 0
    const minLeft = getMinLeft(n, {
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
    for (let right = 0; right < n; right++) {
      if (minLeft[right] === n) continue
      res = Math.max(res, right - minLeft[right] + 1)
    }
    return res
  }
}
