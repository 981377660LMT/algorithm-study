const FAC: number[] = [1]
for (let i = 1; i <= 15; i++) {
  FAC.push(FAC[i - 1] * i)
}

/**
 * 字符串重新排列不含前导0的数字的个数.
 */
function reArrangeDigits(s: string): number {
  const n = s.length
  const counter = Array(10).fill(0)
  for (let i = 0; i < s.length; i++) {
    counter[+s[i]]++
  }

  let res = (n - counter[0]) * FAC[n - 1]
  for (let i = 0; i < 10; i++) {
    res /= FAC[counter[i]]
  }
  return res
}

export { reArrangeDigits }
