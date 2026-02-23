export function estimateTokens(text: string): number {
  // 中文：1 token ≈ 1.5 字符
  // 英文：1 token ≈ 4 字符
  const chineseChars = (text.match(/[\u4e00-\u9fa5]/g) || []).length
  const otherChars = text.length - chineseChars
  return Math.ceil(chineseChars / 1.5 + otherChars / 4)
}

{
  console.log(estimateTokens('sdn你说的那几年·'))
}
