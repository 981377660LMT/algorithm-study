/**
 * @param {string} s
 * @param {string} t
 * @return {boolean}
 * 复杂度O(M+N)
 * @description 缺点：如果t很长 那么都要跑一遍
 */
var isSubsequence1 = function (s: string, t: string): boolean {
  let i = 0
  let j = 0
  while (i < s.length && j < t.length) {
    if (s[i] === t[j]) i++
    j++
  }

  return i === s.length
}

console.log(isSubsequence1('abc', 'ahbgdc'))

// 如果有大量输入的 S，称作 S1, S2, ... , Sk 其中 k >= 10 亿，
// 你需要依次检查它们是否为 T 的子序列。在这种情况下，你会怎样改变代码？
// 像 KMP 算法一样，先用一些时间将长字符串中的数据 提取出来，磨刀不误砍柴功

// 哈希+二分
// 复杂度MlogN
const isSubsequence2 = (s: string, t: string) => {
  const map = new Map<string, number[]>()
  for (let i = 0; i < t.length; i++) {
    !map.has(t[i]) && map.set(t[i], [])
    map.get(t[i])!.push(i)
  }

  console.log(map)
  let index = -Infinity
  // 保证每次index都要比上一次大
  for (const char of s) {
    if (!map.has(char)) return false
    const indexes = map.get(char)!
    // 字母s出现的索引 用二分法找到其中大于index的第一个(bisectRight)
    // 如果l===indexes.length 那么就不存在
    let l = 0
    let r = indexes.length - 1
    while (l <= r) {
      const mid = (l + r) >> 1
      if (indexes[mid] < index) l = mid + 1
      else if (indexes[mid] > index) r = mid - 1
      else l++
    }
    if (l === indexes.length) return false
    index = indexes[l]
  }

  return true
}
console.log(isSubsequence2('acb', 'ahbasdfghytrewgdc'))

type NextPosition = number
// 预处理：记录从该位置开始往后每一个字符第一次出现的位置 复杂度 为s的个数*s的长度
const isSubsequence3 = (s: string, t: string) => {
  const n = t.length
  const chars = new Set(t)
  const maps = Array.from<number, Map<string, NextPosition>>({ length: n + 1 }, () => new Map())
  // 初始化 n表示不存在
  for (const char of chars) {
    maps[n].set(char, n)
  }

  for (let i = n - 1; i >= 0; i--) {
    for (const char of t) {
      if (char === t[i]) maps[i].set(char, i)
      else maps[i].set(char, maps[i + 1].get(char)!)
    }
  }

  console.log(maps)
  // 从第0个查起
  let index = 0
  for (const char of s) {
    console.log(index, char)
    if (!maps[index].has(char) || maps[index].get(char) === n) return false
    index = maps[index].get(char)! + 1
  }
  return true
}

console.log(isSubsequence3('axc', 'ahbgdc'))
export default 1
