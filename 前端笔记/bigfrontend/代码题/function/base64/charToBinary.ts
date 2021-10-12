// 字符转二进制
function charToBinary(text: string) {
  const stringBuilder: string[] = []
  for (const char of text) {
    // 1 bytes = 8bit，将 num 不足8位的0补上
    const num = char.codePointAt(0)!.toString(2).padStart(8, '0')
    stringBuilder.push(num)
  }
  return stringBuilder.join('')
}

if (require.main === module) {
  console.log(charToBinary('as')) // '0110000101110011'
}

export { charToBinary }
