/**
 * @param {any} data
 * @return {string}
 * JSON.stringify() 有额外两个参数，这里并不需要支持。
 */
function stringify(data: any): string {
  // your code here
  if (typeof data === 'bigint')
    throw new Error('Do not know how to serialize a BigInt at JSON.stringify')

  if (typeof data === 'string') return `${data}`

  if (typeof data === 'function') return 'undefined'

  if (data !== data) return 'null'

  if (data === Infinity) return 'null'

  if (data === -Infinity) return 'null'

  if (typeof data === 'number') return `${data}`

  if (typeof data === 'boolean') return `${data}`

  if (data === null) return 'null'

  if (data === undefined) return 'null'

  if (typeof data === 'symbol') return 'null'

  if (data instanceof Date) return `${data.toISOString()}`

  if (Array.isArray(data)) {
    const arr = data.map(el => stringify(el))
    return `[${arr.join(',')}]`
  }

  if (typeof data === 'object') {
    const arr = Object.entries(data).reduce<string[]>((acc, [key, value]) => {
      if (value === undefined) {
        return acc
      }
      acc.push(`"${key}":${stringify(value)}`)
      return acc
    }, [])
    return `{${arr.join(',')}}`
  }

  return ''
}

// console.log(stringify({ 1: 'a' }))

// console.log(JSON.stringify(1n))
