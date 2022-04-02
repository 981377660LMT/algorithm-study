// import assert from 'assert'

// // #region 0. JSON Type
// type JsonValue = JsonPrimitive | JsonObject | JsonArray
// type JsonPrimitive = string | number | boolean | null
// type JsonObject = { [Key in string]?: JsonValue }
// type JsonArray = JsonValue[]
// // #endregion

// // #region 1. normalize string
// // key__snake_bar => KeySnakeBar
// // keyFoo => KeyFoo
// function normalizeKey(key: string): string {
//   if (key.includes('_')) return snakeToPascal(key)
//   return camelToPascal(key)
// }

// function snakeToPascal(snake: string): string {
//   return snake
//     .split(/\_+/) // multi underline
//     .map(capitalize)
//     .join('')
// }

// function camelToPascal(camel: string): string {
//   return capitalize(camel)
// }

// function capitalize(string: string): string {
//   return string[0].toUpperCase() + string.slice(1)
// }
// // #endregion

// // #region 2. main function
// function pascalizeJSONKey(json: JsonValue): JsonValue {
//   if (!isJSONObject(json)) return json

//   const res: JsonObject = {}
//   const visited = new Set<JsonObject>() // 环检测

//   dfs(res, json)
//   return res

//   /**
//    * @description 遍历rawJson里的键值对，将新的键值对放在cur中
//    */
//   function dfs(cur: JsonObject | JsonArray, rawJson: JsonObject | JsonArray): void {
//     for (const [rawKey, rawValue] of Object.entries(rawJson)) {
//       const newKey = normalizeKey(rawKey)

//       if (isJSONObject(rawValue)) {
//         if (visited.has(rawValue)) throw new Error('JSON 中存在环')
//         visited.add(rawValue)

//         const next: JsonObject = {}
//         cur[newKey] = next
//         dfs(next, rawValue)
//       } else {
//         cur[newKey] = rawValue
//       }
//     }
//   }
// }

// function isJSONObject(json: any): json is JsonObject {
//   return Object.prototype.toString.call(json) === '[object Object]'
// }
// // #endregion

// if (require.main === module) {
//   const json = {
//     foo: 'as',
//     bar_1: [1, 2, { snake_bar: 'ok' }],
//     camelCase: { snake_bar: 'ok' },
//     '1': null,
//   }

//   const expected = {
//     Foo: 'as',
//     Bar1: [1, 2, { SnakeBar: 'ok' }],
//     CamelCase: { SnakeBar: 'ok' },
//     '1': null,
//   }

//   assert.deepStrictEqual(pascalizeJSONKey(json), expected)
// }

// export {}

// export {}
