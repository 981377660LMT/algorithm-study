/**
 * @param {number[]} gas
 * @param {number[]} cost
 * @return {number}
 * @description 
 * 在一条环路上有 N 个加油站，其中第 i 个加油站有汽油 gas[i] 升。
   你有一辆油箱容量无限的的汽车，从第 i 个加油站开往第 i+1 个加油站需要消耗汽油 cost[i] 升。你从其中的一个加油站出发，开始时油箱为空。
   如果你可以绕环路行驶一周，则返回出发时加油站的编号，否则返回 -1。
   @summary
   车能开完全程需要满足两个条件：
   车从i站能开到i+1。
   所有站里的油总量要>=车子的总耗油量。
   假设从编号为0站开始，一直到k站都正常，在开往k+1站时车子没油了。这时，应该将起点设置为k+1站。
 */
const canCompleteCircuit = function (gas: number[], cost: number[]): number {
  const sum = (arr: number[]) => arr.reduce((pre, cur) => pre + cur, 0)
  if (sum(gas) < sum(cost)) return -1

  let curGas = 0
  let res = 0
  for (let i = 0; i < gas.length; i++) {
    curGas += gas[i] - cost[i]
    if (curGas < 0) {
      curGas = 0
      res = i + 1
    }
  }

  return res
}

console.log(canCompleteCircuit([1, 2, 3, 4, 5], [3, 4, 5, 1, 2]))
