import { quicktype, InputData, jsonInputForTargetLanguage } from 'quicktype-core'

type JsonValue = string | number | boolean | null | JsonObject | JsonArray
type JsonArray = Array<JsonValue>
interface JsonObject {
  [key: string]: JsonValue
}

export interface ICompressResult {
  /** 压缩后的值（JSON 字符串）*/
  value?: string
  /** 类型定义字符串 */
  valueType?: string
}

/** 最大允许的 Token 长度，压缩后超过此长度会被丢弃. */
const MAX_ALLOWED_TOKENS = 500

/**
 * 压缩 JSON 数据，生成类型定义.
 */
export async function compress(json: JsonValue, typeName: string): Promise<ICompressResult> {
  const result: ICompressResult = {}

  try {
    const compressedJsonStr = JSON.stringify(compressJson(json))
    if (compressedJsonStr !== undefined) {
      const jsonTokens = estimateTokens(compressedJsonStr)
      if (jsonTokens <= MAX_ALLOWED_TOKENS) {
        result.value = compressedJsonStr
      }

      const valueType = await generateQuickType(compressedJsonStr, typeName)
      if (valueType) {
        const typeTokens = estimateTokens(valueType)
        if (typeTokens <= MAX_ALLOWED_TOKENS) {
          result.valueType = valueType
        }
      }
    }
  } catch (error) {
    console.error(error)
    const fallback = JSON.stringify(json)
    if (fallback !== undefined && estimateTokens(fallback) <= MAX_ALLOWED_TOKENS) {
      result.value = fallback
    }
  }

  return result
}

function compressJson(json: JsonValue): JsonValue {
  if (typeof json === 'string') {
    return json.length > 10 ? json.slice(0, 10) + '..' : json
  }
  if (json === null || typeof json !== 'object') {
    return json
  }
  if (Array.isArray(json)) {
    return json.length > 0 ? [compressJson(json[0])] : []
  }
  const res = Object.create(null)
  const keys = Object.keys(json)
  for (const key of keys) {
    res[key] = compressJson((json as any)[key])
  }
  return res
}

/**
 * 使用 quicktype-core 生成类型定义字符串.
 */
async function generateQuickType(jsonString: string, typeName: string): Promise<string> {
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

  return rawOutput
    .replace(/\/\*[\s\S]*?\*\/|\/\/.*/g, '')
    .replace(/export\s+/g, '')
    .replace(/\s+/g, '')
    .replace(/interface/g, 'type ')
    .replace(/;/g, ',')
    .replace(/,}/g, '}')
    .trim()
}

function estimateTokens(text: string): number {
  const chineseChars = (text.match(/[\u4e00-\u9fa5]/g) || []).length
  const otherChars = text.length - chineseChars
  return Math.ceil(chineseChars / 1.5 + otherChars / 4)
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

  // // 1. 简化 JSON (只保留数组第一项，减少类型推断的噪音)
  // const simplified = compressJson(complexResponse)
  // console.log('--- Simplified JSON ---')
  // console.log(JSON.stringify(simplified, null))

  // // 2. 生成 TypeScript 类型定义
  // console.log('\n--- Generated TypeScript Interfaces ---')
  // const tsInterfaces = await generateQuickType(JSON.stringify(simplified), 'foo')
  // console.log(tsInterfaces)

  console.log(await compress(complexResponse, 'Foo'))
})()
