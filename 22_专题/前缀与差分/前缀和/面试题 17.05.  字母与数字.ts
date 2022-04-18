/**
 * @param {string[]} array
 * @return {string[]}
 * 找到最长的子数组，且包含的字母和数字的个数相同。
 * 若存在多个最长子数组，返回左端点下标值最小的子数组
 * 若不存在这样的数组，返回一个空数组。
 * @summary
 * 前缀和数组
 * 当前缀和为0时，假设其下标为n，则说明区间[0, n]内所有元素的和为0，区间长度为n + 1。
 */
const findLongestSubarray = function (array: string[]): string[] {
  // const isNumeric = (x: string) => !isNaN(parseFloat(x)) && isFinite(parseFloat(x))
  const isNumeric = (x: string) => x.codePointAt(0)! >= 48 && x.codePointAt(0)! <= 57

  // 存储前缀和与其第一次出现的位置
  const pre = new Map<number, number>([[0, -1]])
  let sum = 0
  let r = -1
  let maxLen = -1
  let sumWhenMaxLen = -1

  for (let i = 0; i < array.length; i++) {
    const cur = array[i]
    sum += isNumeric(cur) ? 1 : -1
    if (!pre.has(sum)) {
      pre.set(sum, i)
    } else {
      const preIndex = pre.get(sum)!
      const curLen = i - preIndex
      if (curLen > maxLen) {
        maxLen = curLen
        sumWhenMaxLen = sum
        r = i
      }
    }
  }

  console.log(pre, r)
  // console.log(nums, 666, pre) // 有问题
  return array.slice(pre.get(sumWhenMaxLen)! + 1, r + 1)
}

console.log(
  findLongestSubarray([
    'A',
    '1',
    'B',
    'C',
    'D',
    '2',
    '3',
    '4',
    'E',
    '5',
    'F',
    'G',
    '6',
    '7',
    'H',
    'I',
    'J',
    'K',
    'L',
    'M',
  ])
)
console.log(findLongestSubarray(['A', 'A']))
export default 1
