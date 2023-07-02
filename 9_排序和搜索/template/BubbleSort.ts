// BubbleSort 冒泡排序

type Mutable<T> = { -readonly [P in keyof T]: T[P] }

function bubbleSort(arr: Mutable<ArrayLike<number>>): void {
  if (arr.length <= 1) return

  const n = arr.length
  for (let i = 0; i < n - 1; i++) {
    let isSorted = true
    for (let j = 0; j < n - 1 - i; j++) {
      if (arr[j] > arr[j + 1]) {
        const tmp = arr[j]
        arr[j] = arr[j + 1]
        arr[j + 1] = tmp
        isSorted = false
      }
    }
    if (isSorted) break
  }
}

export {}

if (require.main === module) {
  const arr = [1, 4, 2, 5, 3, 6, 7]
  bubbleSort(arr)
  console.log(arr)
}
