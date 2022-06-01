// 输入一个非负整数数组，把数组里所有数字拼接起来排成一个数，打印能拼接出的所有数字中最小的一个。
function minNumber(nums: number[]): string {
  const strs = nums.map(String)
  return (
    strs
      // @ts-ignore
      .sort((a, b) => a + b - (b + a))
      .join('')
  )
}

console.log(minNumber([10, 2]))
console.log(minNumber([3, 30, 34, 5, 9]))
// 输出: "3033459"
// 排序对比不能简单对比字符串,需要连结起来对比
// 如果a+b>b+a则排序后b要在前面
