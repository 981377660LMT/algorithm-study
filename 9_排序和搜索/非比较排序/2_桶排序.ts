// 分桶，排序，合并
// 桶排序是计数排序的升级版。它也是利用函数的映射关系。
const bucketSort = (arr: number[], size: number): number[] => {
  const res: number[] = []
  const min = Math.min.apply(null, arr)
  const max = Math.max.apply(null, arr)
  const count = Math.floor((max - min) / size) + 1
  const bucket = Array.from<number, number[]>({ length: count }, () => [])

  arr.forEach(num => {
    const index = Math.floor((num - min) / size) // 偏移量

    bucket[index].push(num)
  })

  // v8的sort:当数组长度小于等于10的时候，采用插入排序，大于10的时候，采用快排。
  bucket.forEach(bk => res.push(...bk.sort()))

  return res
}

console.log(bucketSort([3, 8, 6, 1, 5, 7, 9, 2, 4], 3))

export {}
