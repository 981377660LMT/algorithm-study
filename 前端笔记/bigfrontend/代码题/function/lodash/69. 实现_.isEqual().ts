/**
 * @param {any} a
 * @param {any} b
 * @return {boolean}
 * 你只需要支持
   基础数据类型
   简单 objects (object literals)
   数组
 * Objects在比较时只需要考虑可枚举属性，且不用考虑prototype中的属性。
   @summary
   1.Use visitedSet to detect circular object 
   2.unwrap until we can compare primitives
 */

function isEqual(a: any, b: any, visited = new Set()): boolean {
  // your code here
  if (a === b) return true

  if (visited.has(a)) return true
  visited.add(a)

  if (typeof a === 'object' && typeof b === 'object') {
    const keysA = [...Object.getOwnPropertySymbols(a), ...Object.keys(a)].sort()
    const keysB = [...Object.getOwnPropertySymbols(b), ...Object.keys(b)].sort()
    if (keysA.length !== keysB.length) return false

    for (let i = 0; i < keysA.length; i++) {
      if (!isEqual(a[keysA[i]], b[keysB[i]], visited)) return false
    }

    return true
  }

  return false
}

if (require.main === module) {
  console.log(isEqual({ a: '1', b: '2' }, { b: '2', a: '1' }))
}
