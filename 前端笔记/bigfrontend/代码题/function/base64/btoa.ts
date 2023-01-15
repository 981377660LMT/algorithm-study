/* eslint-disable @typescript-eslint/no-non-null-assertion */

// 字符串编码为 base64
// 1. 将字符串转换为二进制
// 2. 将二进制6位一组解析，不足6位的用0补齐
// 3. 将 6 位二进制转换为 10 进制，得到 0-63 的数字，对应 BASE64_DIGITS 中的字符
// 4. 最后长度不足 4 的倍数的，用 = 补齐

const BASE64_DIGITS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'

function myBtoa(str: string): string {
  const bin: string[] = []
  for (const char of str) {
    const binStr = char.codePointAt(0)!.toString(2)
    bin.push(binStr.padStart(8, '0'))
  }

  const binStr = bin.join('')
  const res: string[] = []
  for (let i = 0; i < binStr.length; i += 6) {
    const group = binStr.slice(i, i + 6).padEnd(6, '0')
    const index = parseInt(group, 2)
    res.push(BASE64_DIGITS[index])
  }

  const padding = res.length % 4
  if (padding) res.push('='.repeat(4 - padding))
  return res.join('')
}

/**
 * Transform string to base64 in nodejs.
 */
function btoa(str: string): string {
  return Buffer.from(str).toString('base64')
}

if (require.main === module) {
  const str = 'hello world'
  console.log(myBtoa(str))
  console.log(btoa(str))
}

export {}
