// 1 <= strs.length <= 104
// 0 <= strs[i].length <= 100
// 问题:用什么作为map的key? 排序后的str或者26个字母的01值作为key
// 也可以使用 26个质数，代表26个字符，质数相乘结果一定不相等
const groupAnagrams = (strs: string[]): string[][] => {
  const map = new Map<string, string[]>()
  strs.forEach(str => {
    const key = str.split('').sort().join('')
    if (map.has(key)) {
      map.get(key)!.push(str)
    } else {
      map.set(key, [str])
    }
  })
  return [...map.values()]
}

console.log(groupAnagrams(['eat', 'tea', 'tan', 'ate', 'nat', 'bat']))
// 输出: [["bat"],["nat","tan"],["ate","eat","tea"]]
export {}
