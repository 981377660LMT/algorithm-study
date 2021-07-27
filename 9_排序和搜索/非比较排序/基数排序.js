/**
 * name: 基数排序
 * @param  array 待排序数组
 * @param  max 最大位数
 * @description 基数排序是按照低位先排序，然后收集；再按照高位排序，然后再收集；依次类推，直到最高位。
 * 例如，对数字进行排序
 * 1) 创建10个桶，根据个位数依次放入对应的桶中;
 * 2) 根据序号依次将桶内数据放回数组中；
 * 3) 根据十位数、百位数等重复上述步骤，最后一次放回数组时序列有序。
 */
const radixSort = (array, max) => {
  console.time('计数排序耗时')
  const buckets = []
  let unit = 10,
    base = 1
  for (let i = 0; i < max; i++, base *= 10, unit *= 10) {
    for (let j = 0; j < array.length; j++) {
      let index = ~~((array[j] % unit) / base) //依次过滤出个位，十位等等数字
      if (buckets[index] == null) {
        buckets[index] = [] //初始化桶
      }
      buckets[index].push(array[j]) //往不同桶里添加数据
    }
    let pos = 0,
      value
    for (let j = 0, length = buckets.length; j < length; j++) {
      if (buckets[j] != null) {
        while ((value = buckets[j].shift()) != null) {
          array[pos++] = value //将不同桶里数据挨个捞出来，为下一轮高位排序做准备，由于靠近桶底的元素排名靠前，因此从桶底先捞
        }
      }
    }
  }
  console.timeEnd('计数排序耗时')
  return array
}

const array = [3, 44, 38, 5, 47, 15, 36, 26, 27, 2, 46, 4, 19, 50, 48]
console.log('原始array:', array)
const newArr = radixSort(array, 2)
console.log('newArr:', newArr)
