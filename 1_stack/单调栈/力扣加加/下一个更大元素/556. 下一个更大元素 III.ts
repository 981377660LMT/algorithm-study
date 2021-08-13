/**
 * @param {number} n
 * @return {number}
 * 出符合条件的最小整数，其由重新排列 n 中存在的每位数字组成，并且其值大于 n
 * 从尾巴逐位向前找，找到数字比尾巴数字小的第一个位置,与栈顶互换位置
 * 尾巴部分从小到大排序，新的数字即是最小数字
 * 压入栈的尾巴有个特点，一定是升序的。
 */
const nextGreaterElement = function (n: number): number {
  let res = -1
  let arr = n.toString().split('')
  const stack: number[] = []

  for (let i = arr.length - 1; i >= 0; i--) {
    for (const j of stack) {
      // console.log(j, stack)
      if (arr[j] > arr[i]) {
        ;[arr[j], arr[i]] = [arr[i], arr[j]]
        console.log(arr, i)
        const num = parseInt(
          arr.slice(0, i + 1).join('') +
            arr
              .slice(i + 1)
              .sort()
              .join('')
        )
        return num <= 2 ** 31 - 1 ? num : -1
      }
    }
    stack.push(i)
  }
  return res
}

console.log(nextGreaterElement(12))
// 21
console.log(nextGreaterElement(101))
// -1
console.log(nextGreaterElement(2147483486))
