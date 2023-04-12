// 相同的对象
// https://leetcode.com/problems/json-deep-equal/discuss/3406797/JS-Solution-(ES6)
function areDeeplyEqual(o1: any, o2: any): boolean {
  if (o1 === o2) return true

  // 有一个为 null 或者不是对象，直接返回 false
  if (o1 == null || o2 == null) return false
  if (typeof o1 !== 'object' || typeof o2 !== 'object') return false

  const isArray1 = Array.isArray(o1)
  const isArray2 = Array.isArray(o2)
  if (isArray1 !== isArray2) return false

  // 两个都是数组
  if (isArray1 && isArray2) {
    if (o1.length !== o2.length) return false
    for (let i = 0; i < o1.length; i++) {
      if (!areDeeplyEqual(o1[i], o2[i])) return false
    }
    return true
  }

  // 两个都是普通对象
  const keys1 = Object.keys(o1)
  const keys2 = Object.keys(o2)
  if (keys1.length !== keys2.length) return false
  for (const key of keys1) {
    if (!areDeeplyEqual(o1[key], o2[key])) return false
  }
  return true
}

export {}
