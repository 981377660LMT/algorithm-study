/**
 * name: 基数排序
 * @param  nums 待排序数组
 * @param  max 最大位数
 * @description 基数排序是按照低位先排序，然后收集；再按照高位排序，然后再收集；依次类推，直到最高位。
 * 例如，对数字进行排序
   把第 i - 1 轮的结果，根据 nums 中右数第 i 位数，依次放入桶中
   把每个桶的结果连接起来，成为第 i 轮的结果

   时间复杂度O(Row*Col) 这里 Row 是最高位的位数，Col 是数组的长度
 */
function radixSort(nums: readonly number[]): readonly number[] {
  nums = nums.slice()
  const max = Math.max(...nums.map(num => num.toString().length)) // 最大位数

  for (let col = 0; col < max; col++) {
    const buckets: number[][] = Array.from({ length: 10 }, () => [])
    for (const num of nums) {
      const digit = getDigitValue(num, col)
      buckets[digit].push(num)
    }

    nums = buckets.flat()
  }

  return nums
}

// 2进制时
// offset 变成 << 操作
// 模10 变成 &1 操作
function getDigitValue(num: number, digit: number) {
  const offset = 10 ** digit
  return ~~(num / offset) % 10
}

if (require.main === module) {
  const array = [3, 44, 38, 5, 47, 15, 36, 26, 27, 2, 46, 4, 19, 50, 48]
  console.log('原始array:', array)
  const newArr = radixSort(array)
  console.log('newArr:', newArr)
}

export { getDigitValue }
