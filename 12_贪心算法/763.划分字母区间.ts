const BASE = 97

/**
 * @param {string} s
 * @return {number[]}
 * 这个字符串划分为尽可能多的片段，同一字母最多出现在一个片段中
 * 在遍历的过程中相当于是要找每一个字母的边界，如果找到之前遍历过的所有字母的最远边界，说明这个边界就是分割点了。
 * @summary 
 * 统计每一个字符最后出现的位置
   从头遍历字符，并更新字符的最远出现下标，如果找到字符最远出现位置下标和当前下标相等了，则找到了分割点
 */
const partitionLabels = function (s: string): number[] {
  const len = s.length

  const splitLengths: number[] = []
  const rightmostIndex = Array<number>(26).fill(-1)

  for (let i = 0; i < len; i++) {
    rightmostIndex[s[i].codePointAt(0)! - BASE] = i
  }

  let left = 0
  let splitCand = 0
  for (let right = 0; right < len; right++) {
    splitCand = Math.max(splitCand, rightmostIndex[s[right].codePointAt(0)! - BASE])
    if (right === splitCand) {
      splitLengths.push(splitCand - left + 1)
      left = right + 1
    }
  }

  return splitLengths
}

console.log(partitionLabels('ababcbacadefegdehijhklij'))

export default 1
