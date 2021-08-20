/**
 * @param {number[]} gas
 * @param {number[]} cost
 * @return {number}
 * @description 
 * 在一条环路上有 N 个加油站，其中第 i 个加油站有汽油 gas[i] 升。
   你有一辆油箱容量无限的的汽车，从第 i 个加油站开往第 i+1 个加油站需要消耗汽油 cost[i] 升。你从其中的一个加油站出发，开始时油箱为空。
   如果你可以绕环路行驶一周，则返回出发时加油站的编号，否则返回 -1。
   @summary
   计算油的最小值 那么这是最难的地方 把最难的地方最后做 就行了
 */
const canCompleteCircuit = function (gas: number[], cost: number[]): number {
  const sum = (arr: number[]) => arr.reduce((pre, cur) => pre + cur, 0)
  if (sum(gas) < sum(cost)) return -1

  let curGas = 0
  let min = Infinity
  let minGasIndex = 0
  for (let i = 0; i < gas.length; i++) {
    curGas += gas[i] - cost[i]
    if (curGas < min) {
      min = curGas
      minGasIndex = i
    }
  }
  return (minGasIndex + 1) % gas.length
}

console.log(canCompleteCircuit([1, 2, 3, 4, 5], [3, 4, 5, 1, 2]))

export default 1
