// toPairs({ a: 1, b: 2 }) // [['a', 1], ['b', 2]]
// toPairs([2, 4, 8]) // [[0, 2], [1, 4], [2, 8]]
// toPairs('shy') // [['0', 's'], ['1', 'h'], ['2', 'y']]
// toPairs(new Set(['a', 'b', 'c', 'a'])) // [['a', 'a'], ['b', 'b'], ['c', 'c']]

function toPairs(obj) {
  if (obj[Symbol.iterator] instanceof Function && obj.entries instanceof Function)
    // 数组 集合 map
    return Array.from(obj.entries())
  else return Object.entries(obj) // 对象 字符串
}
// console.log(Object.entries('al'))
