/**
 * @param {number[]} nums
 * @return {string}
 * 几个结论
 * 拼接得到的等长字符串结果更大的话，那么原本的整型的数字拼接结果也一定更大
 * 比如 "210" > "102"，那么一定能得到 210 > 102
 * 但两个不等长的字符串就没有这个结论了， 比如 "2" > "10"，但是 2 < 10。
 * @summary
 * 看到要求两个整数 x,y 如何拼接得到结果更大时，就想到先转字符串，然后比较 x+y 和 y+x。这是经验了。
 */
function largestNumber(nums: number[]): string {
  const arr = nums.map(String)
  arr.sort((a, b) => (a + b > b + a ? -1 : 1))
  return BigInt(arr.join('')).toString() // BigInt去除前导零
  return arr.join('').replace(/^0*/, '') || '0'
}

console.log(largestNumber([3, 30, 34, 5, 9]))
// 输出："9534330"

export {}
