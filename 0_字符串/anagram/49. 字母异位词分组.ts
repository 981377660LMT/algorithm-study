// anagram分组
// 1 <= strs.length <= 104
// 0 <= strs[i].length <= 100
// 问题:用什么作为map的key? 排序后的str或者26个字母的01值作为key
// 也可以使用 26个质数，代表26个字符，质数相乘结果一定不相等
const groupAnagrams = (strs: string[]): string[][] => {
  const counter = new Map<string, string[]>()

  strs.forEach(str => {
    const key = str.split('').sort().join('')
    !counter.has(key) && counter.set(key, [])
    counter.get(key)!.push(str)
  })

  return [...counter.values()]
}

console.log(groupAnagrams(['eat', 'tea', 'tan', 'ate', 'nat', 'bat']))
// 输出: [["bat"],["nat","tan"],["ate","eat","tea"]]
export {}
