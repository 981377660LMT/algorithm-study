// 1 <= strs.length <= 104
// 0 <= strs[i].length <= 100
// 问题:用什么作为map的key? 排序后的str或者26个字母的01值作为key
const groupAnagrams = (strs: string[]): string[][] => {
  const map = new Map<string, string[]>()
  strs.forEach(str => {
    const key = str.split('').sort().join('')
    if (map.has(key)) {
      map.set(key, map.get(key)!.concat([str]))
    } else {
      map.set(key, [str])
    }
  })
  return [...map.values()]
}

console.log(groupAnagrams(['eat', 'tea', 'tan', 'ate', 'nat', 'bat']))
// 输出: [["bat"],["nat","tan"],["ate","eat","tea"]]
export {}
