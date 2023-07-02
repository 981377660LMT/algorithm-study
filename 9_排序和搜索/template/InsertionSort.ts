// 插入排序：当数组接近有序时，插入排序的效率很高，近似 O(n)
// 时间复杂度为 O(n^2)，空间复杂度为 O(1)，是稳定排序算法

type Mutable<T> = { -readonly [P in keyof T]: T[P] }

function insertionSort<T>(
  arr: Mutable<ArrayLike<T>>,
  compareFn: (a: T, b: T) => number,
  start = 0,
  end = arr.length
): void {
  if (start < 0) start = 0
  if (end > arr.length) end = arr.length
  if (start >= end) return

  for (let i = start + 1; i < end; i++) {
    for (let j = i; j > start && compareFn(arr[j - 1], arr[j]) > 0; j--) {
      const tmp = arr[j]
      arr[j] = arr[j - 1]
      arr[j - 1] = tmp
    }
  }
}

export {}

if (require.main === module) {
  const arr = [1, 4, 2, 5, 3, 6, 7]
  insertionSort(arr, (a, b) => a - b)
  console.log(arr)
}
