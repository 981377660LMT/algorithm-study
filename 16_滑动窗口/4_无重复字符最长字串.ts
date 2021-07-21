// 思路：
// 找出所有不包含重复字符的字串，返回长度最大的字串
// 滑动窗口
const getLongestSubstring = (str: string): number => {
  let leftPoint = 0
  let maxLength = 0
  let tmpMap = new Map<string, number>()

  for (let rightPoint = 0; rightPoint < str.length; rightPoint++) {
    const element = str[rightPoint]
    // 遇到滑动窗口内的重复字符，看重复字符是谁，如果在滑动窗口里则左指针移动到重复字符的下一位
    // 注意重复值需要在滑动窗口里
    if (tmpMap.has(element) && tmpMap.get(element)! >= leftPoint) {
      leftPoint = tmpMap.get(element)! + 1
    }

    tmpMap.set(element, rightPoint)
    maxLength = Math.max(maxLength, rightPoint - leftPoint + 1)
  }

  return maxLength
}

console.log(getLongestSubstring('abbcdea'))
export {}
