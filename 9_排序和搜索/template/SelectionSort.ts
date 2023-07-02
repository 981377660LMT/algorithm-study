// 每一趟从待排序的数据元素中选择最小（或最大）的一个元素作为首元素，直到所有元素排完为止
const selectSort = (arr: number[]): void => {
  if (arr.length <= 1) return

  for (let i = 0; i < arr.length - 1; i++) {
    let minIndex = i
    for (let j = i; j < arr.length; j++) {
      if (arr[j] < arr[minIndex]) {
        minIndex = j
      }
    }

    const tmp = arr[minIndex]
    arr[minIndex] = arr[i]
    arr[i] = tmp
  }
}

export {}

if (require.main === module) {
  const arr = [4, 1, 2, 5, 3, 6, 7]
  selectSort(arr)
  console.log(arr)
}
