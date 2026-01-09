export function safeParse<T>(str: unknown, defaultValue: T): T {
  if (typeof str !== 'string') return defaultValue
  try {
    const result = JSON.parse(str)
    // JSON.parse 不会返回 undefined，但为了类型安全
    return result == null ? defaultValue : result
  } catch {
    return defaultValue
  }
}

function ensureArray(val: unknown) {
  try {
    let parsed = val
    // 循环解析，处理多层序列化的情况
    while (typeof parsed === 'string') {
      parsed = JSON.parse(parsed || '[]')
    }
    return Array.isArray(parsed) ? parsed : []
  } catch (e) {
    return []
  }
}
