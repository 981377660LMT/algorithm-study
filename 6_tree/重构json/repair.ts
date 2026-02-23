import { jsonrepair } from 'jsonrepair'

try {
  const json = `{
  "exceptionIds": ['{{urlExceptionId.value}}'],
  "pageNum": 1,
  "pageSize": 1
}`

  const repaired = jsonrepair(json)

  console.log(repaired) // '{"name": "John"}'
  console.log(JSON.parse(repaired))
} catch (err) {
  console.error(err)
}

// try {
//   JSON.parse(str)
// } catch (e) {
//   console.log(e.message)
// }

function safeParseRecordString<T = Record<string, string>>(
  value: unknown,
  fallback?: T
): T | undefined {
  const isObject = (v: unknown) => v !== null && typeof v === 'object' && !Array.isArray(v)

  if (isObject(value)) return value as T

  if (typeof value === 'string') {
    try {
      const parsed = JSON.parse(value)
      if (isObject(parsed)) return parsed as T
    } catch {}

    try {
      // 将裸露的 {{...}} 包裹上双引号
      let fixedValue = value
        .replace(/:\s*(\{\{[^}]+\}\})(?=[,}])/g, ':"$1"')
        .replace(/\[\s*(\{\{[^}]+\}\})\s*\]/g, '["$1"]')
      const repaired = jsonrepair(fixedValue)
      const parsed = JSON.parse(repaired)
      if (isObject(parsed)) return parsed as T
    } catch (e) {
      console.warn('safeParseRecordString failed:', e)
    }
  }

  return fallback
}
