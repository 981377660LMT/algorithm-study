/**
 * @param {number} m
 * @param {number} n
 * @param {number} k
 * @return {number}
 * 地上有一个m行n列的方格，从坐标 [0,0] 到坐标 [m-1,n-1]
 * 它每次可以向左、右、上、下移动一格（不能移动到方格外），也不能进入行坐标和列坐标的数位之和大于k的格子
 * @summary 因为问范围 所以只要考虑向下行走和向右行走
 */
const movingCount = function (m: number, n: number, k: number): number {}

console.log(movingCount(2, 3, 1))

// 数位和
function sum(num: number) {
  let res = 0
  while (num) {
    res += num % 10
    res = ~~(num / 10)
  }
  return res
}

// x数位和为sum 求(x+1)数位和
function sumDiff(sum: number) {
  return (sum + 1) % 10 ? sum + 1 : sum - 8
}

export {}
