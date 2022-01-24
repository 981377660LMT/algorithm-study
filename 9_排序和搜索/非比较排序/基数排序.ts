/**
 * name: 基数排序
 * @param  array 待排序数组
 * @param  max 最大位数
 * @description 基数排序是按照低位先排序，然后收集；再按照高位排序，然后再收集；依次类推，直到最高位。
 * 例如，对数字进行排序
 * 1) 创建10个桶，根据个位数依次放入对应的桶中;
 * 2) 根据序号依次将桶内数据放回数组中；
 * 3) 根据十位数、百位数等重复上述步骤，最后一次放回数组时序列有序。
 * 时间复杂度(10n)
 */
function radixSort(array: readonly number[]): readonly number[] {
  const maxDigitLength = Math.max(...array.map(num => num.toString().length))

  for (let curDigit = 0; curDigit < maxDigitLength; curDigit++) {
    const bucket = Array.from<unknown, number[]>({ length: 10 }, () => [])
    for (const num of array) {
      // 求出每一个元素的个、十、百位的值
      const index = getDigitValue(num, curDigit)
      bucket[index].push(num)
    }

    array = bucket.flat()
  }

  return array
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
