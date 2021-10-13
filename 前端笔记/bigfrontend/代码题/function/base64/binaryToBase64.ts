import { charToBinary } from './charToBinary'

const store = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'

// 将二进制数据每 6bit 位替换成一个 base64 字符
// 1 bytes = 8bit，6bit 位替换成一个 base64 字符
// 所以每 3 bytes 的数据，能成功替换成 4 个 base64 字符
function binaryTobase64(binary: string): string {
  const [normalizedBinary, suffix] = normalize(binary)

  const stringBuilder: string[] = []
  // 按 6bit 一组转换
  for (let i = 0; i < normalizedBinary.length; i += 6) {
    const slice = normalizedBinary.slice(i, i + 6)
    const index = parseInt(slice, 2)
    stringBuilder.push(store[index])
  }

  return stringBuilder.join('') + suffix
}

/**
 *
 * @param binary
 * 二进制长度要凑成6的倍数：
 * 8bit：补4个0
 * 16bit:补2个0
 * 加几个等号其实没啥特殊意义 只是在binary上标记原始信息
 */
function normalize(binary: string): [string, string] {
  if (binary.length % 24 === 8) return [binary + '0000', '==']
  else if (binary.length % 24 === 16) return [binary + '00', '=']
  else return [binary, '']
}

console.log(binaryTobase64(charToBinary('as')))
console.log(Buffer.from('as').toString('base64'))
