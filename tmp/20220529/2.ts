function isNumeric(str: string): boolean {
  if (typeof str !== 'string') return false
  const num = parseFloat(str)
  // 注意不要用Number.isFinite(str)/Number.isNaN(str) 因为会先把字符串转成数字 从而所有字符串都被认为false
  return !isNaN(num) && isFinite(Number(str))
}

function discountPrices(sentence: string, discount: number): string {
  const words = sentence.split(' ')
  const res: string[] = []

  for (const word of words) {
    if (word[0] === '$' && isNumeric(word.slice(1))) {
      const price = parseFloat(word.slice(1))
      const discountedPrice = (price * (100 - discount)) / 100
      res.push(`$${discountedPrice.toFixed(2)}`)
    } else {
      res.push(word)
    }
  }

  return res.join(' ')
}

export {}
console.log(/^\d+$/.test(''))
