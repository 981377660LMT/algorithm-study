/* eslint-disable implicit-arrow-linebreak */

// 求出任意一个最长公共子串(数组) 更新res时记录一下路径
function calLCS(s1: string, s2: string): string {
  const n1 = s1.length
  const n2 = s2.length
  const dp = Array.from<number, number[]>({ length: n1 + 1 }, () => Array(n2 + 1).fill(0))
  let last = [-1, -1]

  for (let i = 1; i < n1 + 1; i++) {
    for (let j = 1; j < n2 + 1; j++) {
      if (s1[i - 1] === s2[j - 1]) {
        const cand = dp[i - 1][j - 1] + 1
        if (cand > dp[i][j]) {
          dp[i][j] = cand
          last = [i, j]
        }
      }
    }
  }

  const sb: string[] = []
  let [i, j] = last
  while (i > 0 && j > 0 && s1[i - 1] === s2[j - 1]) {
    sb.push(s1[i - 1])
    i--
    j--
  }

  return sb.reverse().join('')
}

const str1 = '123456'
const str2 = '1234'
console.log(calLCS(str1, str2))

export {}
