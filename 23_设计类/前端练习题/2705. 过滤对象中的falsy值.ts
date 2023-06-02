// 2705. 过滤对象中的falsy值

type JsonValue = JsonPrimitive | JsonObject | JsonArray
type JsonPrimitive = string | number | boolean | null
type JsonObject = { [Key in string]?: JsonValue }
type JsonArray = JsonValue[]

function compactObject(obj: JsonValue): JsonValue {
  if (!isObj(obj)) return obj
  if (Array.isArray(obj)) return obj.filter(Boolean).map(compactObject) // !先filter再map速度会快一些

  const res: JsonObject = {}
  for (const key of Object.keys(obj)) {
    const child = compactObject(obj[key]!)
    if (child) res[key] = child
  }
  return res
}

function isObj(o: unknown): o is JsonObject | JsonArray {
  return typeof o === 'object' && o !== null
}

export {}
