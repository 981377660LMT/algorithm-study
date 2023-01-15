// 罗马数字转整数
const MAPPING: Record<string, number> = {
  I: 1,
  V: 5,
  X: 10,
  L: 50,
  C: 100,
  D: 500,
  M: 1000
}

/**
 * @param {string} s
 * @return {number}
 * if currInt is smaller than nextInt, it indicates the substraction.
   For example of IV currInt is 1 and nextInt is 5, then -1+5 = 4
   For example of VI currInt is 5 and nextInt is 1, then 5+1= 6
 */
function romanToInt(s: string): number {
  let res = 0

  // 小就减 不小就加
  for (let i = 0; i < s.length; i++) {
    MAPPING[s[i]] < MAPPING[s[i + 1]] ? (res -= MAPPING[s[i]]) : (res += MAPPING[s[i]])
  }

  return res
}

console.log(romanToInt('CXXIII'))
// 123
console.log(romanToInt('MCMXCIX'))
// 1999

export {}
