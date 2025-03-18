// 相同的对象
// https://leetcode.com/problems/json-deep-equal/discuss/3406797/JS-Solution-(ES6)
function areDeeplyEqual(o1: any, o2: any): boolean {
  if (o1 === o2) return true
  if (!isObject(o1) || !isObject(o2)) return false

  if (Array.isArray(o1) !== Array.isArray(o2)) return false

  const keys1 = Object.keys(o1)
  const keys2 = Object.keys(o2)
  if (keys1.length !== keys2.length) return false
  for (const key of keys1) {
    if (!areDeeplyEqual(o1[key], o2[key])) return false
  }
  return true
}

function isObject(obj: unknown): boolean {
  return obj != null && typeof obj === 'object'
}

export { areDeeplyEqual }
