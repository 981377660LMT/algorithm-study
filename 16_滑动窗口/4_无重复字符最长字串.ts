// 思路：
// 找出所有不包含重复字符的字串，返回长度最大的字串
// 滑动窗口
const getLongestSubstring = (str: string): number => {
  let l = 0
  let res = 0
  let map = new Map<string, number>()

  for (let r = 0; r < str.length; r++) {
    const cur = str[r]
    // 遇到滑动窗口内的重复字符，看重复字符是谁，如果在滑动窗口里则左指针移动到重复字符的下一位
    // 注意重复值需要在滑动窗口里
    if (map.has(cur) && map.get(cur)! >= l) {
      l = map.get(cur)! + 1
    }

    map.set(cur, r)
    res = Math.max(res, r - l + 1)
  }

  return res
}

console.log(getLongestSubstring('abbcdea'))
export {}
