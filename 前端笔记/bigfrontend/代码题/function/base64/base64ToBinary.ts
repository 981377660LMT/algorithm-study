const store = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'
const codeToIndex = {} as Record<string, number>
for (let i = 0; i < store.length; i++) {
  codeToIndex[store[i]] = i
}

function base64ToBinary(base64: string): string {
  check(base64)
  const len = base64.length
  const stringBuilder: string[] = []
  const [normalizedBase64, suffix] = normalize(base64)

  // every 3 bytes take 4 b64 chars
  for (let i = 0; i < len; i += 4) {
    const n0 = codeToIndex[normalizedBase64[i]]
    const n1 = codeToIndex[normalizedBase64[i + 1]]
    const n2 = codeToIndex[normalizedBase64[i + 2]]
    const n3 = codeToIndex[normalizedBase64[i + 3]]

    const binaryOfThreeBytes = (n0 << 18) + (n1 << 12) + (n2 << 6) + n3

    stringBuilder.push(
      String.fromCodePoint((binaryOfThreeBytes >>> 16) & 0xff),
      String.fromCodePoint((binaryOfThreeBytes >>> 8) & 0xff),
      String.fromCodePoint(binaryOfThreeBytes & 0xff)
    )
  }

  return (suffix ? stringBuilder.slice(0, -suffix.length) : stringBuilder).join('')
}

function normalize(base64: string): [string, string] {
  const len = base64.length
  const suffix = base64[len - 2] === '=' ? 'AA' : base64[len - 1] === '=' ? 'A' : ''
  const normalizedBase64 = suffix ? `${base64.slice(0, -suffix.length)}${suffix}` : base64
  return [normalizedBase64, suffix]
}

function check(base64: string) {
  if (base64.length % 4 !== 0) throw new Error('Invalid input')
  // "A==="   Expected function to throw an exception.
  if (/={3,}/.test(base64)) throw new Error('Invalid input')
  // "QkZFLmRld=g="   Expected function to throw an exception.
  if (/=+[^=]+=+/.test(base64)) throw new Error('Invalid input')
}

export {}

console.log(base64ToBinary('YXM='))
console.log(Buffer.from('YXM=', 'base64').toString())
