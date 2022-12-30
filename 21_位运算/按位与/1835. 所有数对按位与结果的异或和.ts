// 1835. 所有数对按位与结果的异或和
// 布尔代数里异或就是加，与就是乘，所以符合分配律
// a&c ^ b&c = (a^b) & c
function getXORSum(arr1: number[], arr2: number[]): number {
  let xor1 = arr1.reduce((pre, cur) => pre ^ cur, 0)
  let xor2 = arr2.reduce((pre, cur) => pre ^ cur, 0)
  return xor1 & xor2
}

console.log(getXORSum([1, 2, 3], [6, 5]))
// 输出：0
// 解释：列表 = [1 AND 6, 1 AND 5, 2 AND 6, 2 AND 5, 3 AND 6, 3 AND 5] = [0,1,2,0,2,1] ，
// 异或和 = 0 XOR 1 XOR 2 XOR 0 XOR 2 XOR 1 = 0 。
