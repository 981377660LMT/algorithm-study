// 使用 JSON.parse(JSON.stringify(object)) 的弊端

// 会忽略 undefined
// 会忽略 symbol
// 不能序列化函数
// 不能解决循环引用的对象

interface CacheItem {
  original: Record<any, any>
  copy: Record<any, any>
}

/**
 * Deep copy the given object considering circular structure.
 * This function caches all nested objects and its copies.
 * If it detects circular structure, use cached copy to avoid infinite loop.
 *
 * @param {*} obj
 * @param {Array<Object>} cache
 * @return {*}
 */
const deepCopy = (obj: any, cache: CacheItem[] = []) => {
  if (obj === null || typeof obj !== 'object') return obj

  const hit = cache.filter(v => v.original === obj)[0]
  if (hit) return hit.copy

  const copy = Array.isArray(obj) ? [] : ({} as Record<any, any>)
  cache.push({
    copy,
    original: obj,
  })

  // 如果要拷贝Symbol的话 这里获取Symbol的键
  Object.keys(obj).forEach(key => {
    copy[key] = deepCopy(obj[key], cache)
  })

  return copy
}
