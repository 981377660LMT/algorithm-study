import { writeFileSync } from 'fs'
import path from 'path'

type JsonValue = string | number | boolean | null | JsonObject | JsonArray
interface JsonObject {
  [key: string]: JsonValue
}
interface JsonArray extends Array<JsonValue> {}

/**
 * 递归地对JSON对象或数组的键进行排序，以实现规范化。
 * @param json 要规范化的JSON值。
 * @returns 键已排序的规范化JSON值。
 */
function normalizeJson(json: JsonValue): JsonValue {
  if (json === null || typeof json !== 'object') {
    return json
  }
  if (Array.isArray(json)) {
    return json.map(normalizeJson)
  }
  const sortedKeys = Object.keys(json).sort()
  const res: Record<string, JsonValue> = {}
  for (const key of sortedKeys) {
    res[key] = normalizeJson(json[key])
  }
  return res
}

export {}

if (require.main === module) {
  const obj1 = {}
  // 写入文件
  const componentData = normalizeJson(obj1)

  const outputPath = path.join(__dirname, 'normalized5.json')
  writeFileSync(outputPath, JSON.stringify(componentData, null, 2))
}
