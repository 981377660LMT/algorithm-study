// 使用Uint8Array加速
const maxVowels = (s: string, k: number) => {
  // const weight = new Uint8Array(123)
  // let max = 0
  // weight[97] = weight[101] = weight[105] = weight[111] = weight[117] = 1
  // for (let i = 0; i < k; ++i) {
  //   max += weight[s.codePointAt(i)]
  // }
  // for (let i = 0, cur = max; i < s.length - k; ++i) {
  //   cur += weight[s.codePointAt(i + k)] - weight[s.codePointAt(i)]
  //   cur > max && (max = cur)
  // }
  // return max

  // 使用位运算减少空间复杂度
  const BASE = 97
  const weight = 1065233 // 1 + (1 << 4) + (1 << 8) + (1 << 14) + (1 << 20);
  let max = 0

  for (let i = 0; i < k; ++i) {
    max += (weight >> (s.codePointAt(i)! - BASE)) & 1
  }

  for (let i = 0, cur = max; i < s.length - k; ++i) {
    cur +=
      ((weight >> (s.codePointAt(i + k)! - BASE)) & 1) -
      ((weight >> (s.codePointAt(i)! - BASE)) & 1)
    cur > max && (max = cur)
  }

  return max
}

console.log(maxVowels('leetcode', 3))

export {}
