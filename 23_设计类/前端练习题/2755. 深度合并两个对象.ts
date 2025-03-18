// 递归合并两个JSON对象(对象是JSON.parse()解析的结果)

// 1. 类型不同，返回第二个对象
// 2. 都是基本类型，返回第二个对象
// 3. 都是数组，递归合并
// 4. 都是对象，递归合并

function deepMerge(obj1: any, obj2: any): any {
  if (obj1 === void 0) return obj2
  if (obj2 === void 0) return obj1

  if (getType(obj1) !== getType(obj2)) return obj2
  if (!isObj(obj1) || !isObj(obj2)) return obj2

  // 都是数组
  if (Array.isArray(obj1) && Array.isArray(obj2)) {
    const n1 = obj1.length
    const n2 = obj2.length
    const n = Math.max(n1, n2)
    const res = Array(n)
    for (let i = 0; i < n; i++) {
      res[i] = deepMerge(obj1[i], obj2[i])
    }
    return res
  }

  // 都是对象
  const res: any = { ...obj1 }
  for (const key in obj2) {
    if (Object.prototype.hasOwnProperty.call(obj2, key)) {
      res[key] = deepMerge(obj1[key], obj2[key])
    }
  }
  return res
}

function isObj(o: unknown): o is object {
  return typeof o === 'object' && o !== null
}

function getType(o: unknown): string {
  return Object.prototype.toString.call(o).slice(8, -1)
}

/**
 * let obj1 = {"a": 1, "c": 3}, obj2 = {"a": 2, "b": 2};
 * deepMerge(obj1, obj2); // {"a": 2, "c": 3, "b": 2}
 */

export {}
