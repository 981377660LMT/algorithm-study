// a&c ^ b&c = (a^b) & c
function getXORSum(arr1: number[], arr2: number[]): number {
  let xor1 = 0
  let xor2 = 0
  for (const num of arr1) {
    xor1 ^= num
  }
  for (const num of arr2) {
    xor2 ^= num
  }
  return xor1 & xor2
}

console.log(getXORSum([1, 2, 3], [6, 5]))
// 输出：0
// 解释：列表 = [1 AND 6, 1 AND 5, 2 AND 6, 2 AND 5, 3 AND 6, 3 AND 5] = [0,1,2,0,2,1] ，
// 异或和 = 0 XOR 1 XOR 2 XOR 0 XOR 2 XOR 1 = 0 。
