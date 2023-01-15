// base64 to string
// 1. 验证 base64 是否合法(长度是否为 4 的倍数,末尾最多 2 个字符是 =)
// 2. 去除末尾的 = , 用'A'(0)作为占位符补齐
// !3. 每4个字符一组(3个字符对应4个base64字符),转为uint24,然后转为3个字符
// 4. 末尾去除paddingZero个占位符

const BASE64_DIGITS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'
const BASE64_DIGITS_MAPPING: Record<string, number> = {}
for (let i = 0; i < BASE64_DIGITS.length; i++) {
  BASE64_DIGITS_MAPPING[BASE64_DIGITS[i]] = i
}

// throw error if base64 is not valid
function myAtob(base64: string): string {
  check(base64)
  const n = base64.length
  let paddingZero = ''
  if (base64[n - 2] === '=') paddingZero = 'AA'
  else if (base64[n - 1] === '=') paddingZero = 'A'
  base64 = paddingZero ? `${base64.slice(0, -paddingZero.length)}${paddingZero}` : base64

  const res: string[] = []
  for (let i = 0; i < n; i += 4) {
    const i0 = BASE64_DIGITS_MAPPING[base64[i]]
    const i1 = BASE64_DIGITS_MAPPING[base64[i + 1]]
    const i2 = BASE64_DIGITS_MAPPING[base64[i + 2]]
    const i3 = BASE64_DIGITS_MAPPING[base64[i + 3]]
    const uint24 = (i0 << 18) + (i1 << 12) + (i2 << 6) + i3
    res.push(
      String.fromCodePoint(uint24 >>> 16),
      String.fromCodePoint((uint24 >>> 8) & 0xff),
      String.fromCodePoint(uint24 & 0xff)
    )
  }

  return paddingZero ? res.slice(0, -paddingZero.length).join('') : res.join('')

  function check(base64: string): void {
    if (base64.length % 4) throw new Error('Invalid base64')
    if (/={3,}/.test(base64)) throw new Error('Invalid input') // "A==="   结尾最多有两个等号
    if (/=+[^=]+=+/.test(base64)) throw new Error('Invalid input') // "QkZFLmRld=g="   等号中间不能有非等号字符
  }
}

function atob(base64: string): string {
  return Buffer.from(base64, 'base64').toString()
}

if (require.main === module) {
  const base64 = 'SGVsbG8gV29ybGQh'
  // console.log(myAtob(base64))
  console.log(atob(base64)) // Hello World!
  console.log(atob('5aW9')) // 好
  console.log(myAtob('YQ=='))
  console.log(atob('YQ=='))
}

export {}
