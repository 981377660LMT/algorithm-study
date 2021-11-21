// 对> 2的53次方-1的数，JS为我们提供了内置基本数据类型BigInt
// 能被5整除的数，个位只能是0或5
function prefixesDivBy5(nums: number[]): boolean[] {
  const res: boolean[] = []
  let prefix = 0

  for (let i = 0; i < nums.length; i++) {
    prefix = ((prefix << 1) + nums[i]) % 5
    res.push(prefix === 0)
  }

  return res
}

console.log(prefixesDivBy5([0, 1, 1]))
// 输出：[true,false,false]
// 解释：
// 输入数字为 0, 01, 011；也就是十进制中的 0, 1, 3 。只有第一个数可以被 5 整除，因此 answer[0] 为真。
