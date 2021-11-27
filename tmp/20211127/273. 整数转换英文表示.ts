const radix = [1000000000, 1000000, 1000, 1]
const LESS_THAN_20 = [
  '',
  'One',
  'Two',
  'Three',
  'Four',
  'Five',
  'Six',
  'Seven',
  'Eight',
  'Nine',
  'Ten',
  'Eleven',
  'Twelve',
  'Thirteen',
  'Fourteen',
  'Fifteen',
  'Sixteen',
  'Seventeen',
  'Eighteen',
  'Nineteen',
]

const TENS = [
  '',
  'Ten',
  'Twenty',
  'Thirty',
  'Forty',
  'Fifty',
  'Sixty',
  'Seventy',
  'Eighty',
  'Ninety',
]

const THOUSANDS = ['Billion', 'Million', 'Thousand', '']

function numberToWords(num: number): string {
  if (num === 0) return 'Zero'
  const sb: string[] = []

  for (let i = 0; i < radix.length; i++) {
    const [div, mod] = [~~(num / radix[i]), num % radix[i]]
    if (div === 0) continue
    sb.push(trans(div))
    sb.push(THOUSANDS[i])
    sb.push(' ')
    num = mod
  }

  return sb.join('').trim()

  // 处理1000以下的数
  function trans(num: number): string {
    if (num === 0) return ''
    if (num < 20) return LESS_THAN_20[num] + ' '
    if (num < 100) return TENS[~~(num / 10)] + ' ' + trans(num % 10)
    return LESS_THAN_20[~~(num / 100)] + ' Hundred ' + trans(num % 100)
  }
}

console.log(numberToWords(1234567891))
// 输入：num = 1234567891
// 输出："One Billion Two Hundred Thirty Four Million Five Hundred Sixty Seven Thousand Eight Hundred Ninety One"
