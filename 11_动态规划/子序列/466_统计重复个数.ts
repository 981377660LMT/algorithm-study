/**
 * @param {string} s1
 * @param {number} n1
 * @param {string} s2
 * @param {number} n2
 * @return {number}
 * str = [s, n] 表示 str 由 n 个字符串 s 连接构成
 * 请你找出一个最大整数 m ，以满足 str = [str2, m] 可以从 str1 获得。
 */
const getMaxRepetitions = function (s1: string, n1: number, s2: string, n2: number): number {
  let s1Count = 0
  let s2Count = 0
  let s2p = 0

  // 计算n1个s1中包含了多少个s2
  while (s1Count < n1) {
    for (let index = 0; index < s1.length; index++) {
      const s1letter = s1[index]
      if (s1letter === s2[s2p]) s2p++
      if (s2p === s2.length) {
        s2Count++
        s2p = 0
      }
    }

    s1Count++

    // 找到了循环节
    if (s2p === 0) {
      // 一共需要循环多少次
      //这里计数乘循环的次数，继续循环 因为s1Count还可能是小于n1的，循环节点不能整除
      const times = Math.floor(n1 / s1Count)
      s1Count *= times
      s2Count *= times
    }
  }

  // s2Count：包含多少个s2
  return Math.floor(s2Count / n2)
}
