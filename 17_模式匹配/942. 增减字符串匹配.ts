// 贪心法，[0,N]的数，如果出现D，那么从大的一端取，出现I，从小的一端取
function diStringMatch(s: string): number[] {
  const res: number[] = []
  const nums = Array.from<unknown, number>({ length: s.length + 1 }, (_, i) => i)
  let left = 0
  let right = s.length

  for (let i = 0; res.length <= s.length; i++) {
    if (s[i] === 'I') {
      res.push(nums[left])
      left++
    } else {
      res.push(nums[right])
      right--
    }
  }

  return res
}

console.log(diStringMatch('IDID'))

// 示例 1：

// 输入："IDID"
// 输出：[0,4,1,3,2]
// 示例 2：

// 输入："III"
// 输出：[0,1,2,3]
