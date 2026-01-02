// https://github.com/glideapps/quicktype/issues/2763
// 注意要用版本 quicktype-core@23.0.146，否则浏览器内报错

import { quicktype, InputData, jsonInputForTargetLanguage } from 'quicktype-core'

type JsonValue = string | number | boolean | null | JsonObject | JsonArray
interface JsonObject {
  [key: string]: JsonValue
}
type JsonArray = JsonValue[]

/**
 * 递归遍历 JSON 对象，将所有数组保留第一项。
 */
function keepFirstItemInArrays(json: JsonValue): JsonValue {
  const isObject = json !== null && typeof json === 'object'
  if (!isObject) return json

  if (Array.isArray(json)) {
    if (!json.length) return []
    return [keepFirstItemInArrays(json[0])]
  }

  const res = Object.create(null)
  for (const key in json) {
    if (Object.prototype.hasOwnProperty.call(json, key)) {
      res[key] = keepFirstItemInArrays(json[key])
    }
  }
  return res
}

/**
 * 使用 quicktype-core 生成类型定义字符串.
 */
async function generateQuickType(json: JsonValue, typeName: string): Promise<string> {
  const jsonString = JSON.stringify(json)

  const jsonInput = jsonInputForTargetLanguage('typescript')
  await jsonInput.addSource({
    name: typeName,
    samples: [jsonString]
  })

  const inputData = new InputData()
  inputData.addInput(jsonInput)

  const res = await quicktype({
    inputData,
    lang: 'typescript',
    rendererOptions: {
      'just-types': 'true',
      'explicit-unions': 'true'
    }
  })
  const rawOutput = res.lines.join('\n')

  // 执行压缩：
  return rawOutput
    .replace(/\/\*[\s\S]*?\*\/|\/\/.*/g, '')
    .replace(/export\s+/g, '')
    .replace(/\s+/g, '') // 移除所有空格，极致压缩
    .replace(/interface/g, 'type ')
    .replace(/;/g, ',')
    .replace(/,}/g, '}')
    .trim()
}

;(async () => {
  const complexResponse = {
    userId: 101,
    userName: 'Alice',
    roles: ['admin', 'editor', 'viewer'],
    orders: [
      {
        id: 'ord_1',
        total: 100,
        items: [
          { productId: 'p1', name: 'Book', price: 20 },
          { productId: 'p2', name: 'Pen', price: 5 }
        ]
      },
      {
        id: 'ord_2',
        total: 200,
        items: []
      }
    ],
    meta: {
      tags: [
        { id: 1, label: 'urgent' },
        { id: 2, label: 'home' }
      ]
    }
  }
  // 1. 简化 JSON (只保留数组第一项，减少类型推断的噪音)
  const simplified = keepFirstItemInArrays(complexResponse)
  console.log('--- Simplified JSON ---')
  console.log(JSON.stringify(simplified, null))

  // 2. 生成 TypeScript 类型定义
  console.log('\n--- Generated TypeScript Interfaces ---')
  const tsInterfaces = await generateQuickType(complexResponse, 'foo')
  console.log(tsInterfaces)
})()
