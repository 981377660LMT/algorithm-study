import { quicktype, InputData, jsonInputForTargetLanguage } from 'quicktype-core'

import { JsonValue } from './types'

export interface ICompressOptions {
  /** 压缩后 value 的最大字符数，超过则舍弃 value 只保留 type */
  maxValueLength?: number
}

export interface ICompressResult {
  /** 压缩后的值（JSON 字符串），如果太大则为 undefined */
  value?: string
  /** 类型定义字符串 */
  valueType: string
  /** 压缩策略说明 */
  strategy: 'compressed' | 'original' | 'type-only' | 'error'
}

const DEFAULT_OPTIONS: Required<ICompressOptions> = {
  maxValueLength: 200
}

/**
 * 压缩 JSON 数据，生成类型定义。
 *
 * 1. 将数组保留第一项，减少数据量
 * 2. 生成 TypeScript 类型定义
 * 3. 如果压缩后(value+type)比原始 value 更大，则使用原始值（不生成类型）
 * 4. 如果压缩后 value 仍然太大，只保留 type
 */
export async function compress(
  json: JsonValue,
  typeName: string,
  options?: ICompressOptions
): Promise<ICompressResult> {
  const opts = { ...DEFAULT_OPTIONS, ...options }

  try {
    const originalJsonStr = JSON.stringify(json)
    const originalTotal = originalJsonStr.length

    const compressedValue = retainFirstElement(json)
    const compressedJsonStr = JSON.stringify(compressedValue)

    const valueType = await generateQuickType(compressedJsonStr, typeName)
    if (!valueType) {
      if (originalTotal <= opts.maxValueLength) {
        return { value: originalJsonStr, valueType: '', strategy: 'original' }
      }
      return { value: undefined, valueType: '', strategy: 'type-only' }
    }

    const compressedTotal = compressedJsonStr.length + valueType.length
    if (compressedTotal >= originalTotal) {
      if (originalTotal <= opts.maxValueLength) {
        return {
          value: originalJsonStr,
          valueType: '',
          strategy: 'original'
        }
      }
    }

    if (compressedJsonStr.length <= opts.maxValueLength) {
      return {
        value: compressedJsonStr,
        valueType,
        strategy: 'compressed'
      }
    }

    return {
      value: undefined,
      valueType,
      strategy: 'type-only'
    }
  } catch (error) {
    console.error('compress failed', error)
    return {
      value: JSON.stringify(json),
      valueType: '',
      strategy: 'error'
    }
  }
}

/**
 * 压缩策略：
 * 1. 数组仅保留首项
 * 2. 字符串截断至 20 字符
 * 3. 对象属性最多保留 10 个
 */
function retainFirstElement(json: JsonValue): JsonValue {
  if (typeof json === 'string') {
    return json.length > 20 ? json.slice(0, 20) + '..' : json
  }
  if (json === null || typeof json !== 'object') {
    return json
  }
  if (Array.isArray(json)) {
    return json.length > 0 ? [retainFirstElement(json[0])] : []
  }

  const res = Object.create(null)
  const keys = Object.keys(json)
  const limitedKeys = keys.slice(0, 10)

  for (const key of limitedKeys) {
    res[key] = retainFirstElement((json as any)[key])
  }

  if (keys.length > 10) {
    res['__more__'] = '...'
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
