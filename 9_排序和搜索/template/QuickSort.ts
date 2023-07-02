// QucikSort 快速排序三路快排
// 不稳定

type Mutable<T> = { -readonly [P in keyof T]: T[P] }

function quickSort<T>(
  arr: Mutable<ArrayLike<T>>,
  compareFn: (a: T, b: T) => number,
  start = 0,
  end = arr.length
): void {
  if (start < 0) start = 0
  if (end > arr.length) end = arr.length
  if (start >= end) return

  partition(start, end - 1)

  function partition(left: number, right: number): void {
    if (left >= right) return

    // 优化，随机取标定点，以解决近乎有序的列表
    const randIndex = randint(left, right)
    swap(left, randIndex)

    const pivot = arr[left]
    let leftPoint = left
    let midPoint = left
    let rightPoint = right
    while (midPoint <= rightPoint) {
      const compared = compareFn(arr[midPoint], pivot)
      if (compared < 0) {
        swap(leftPoint, midPoint)
        leftPoint++
        midPoint++
      } else if (compared > 0) {
        swap(midPoint, rightPoint)
        rightPoint--
      } else {
        midPoint++
      }
    }

    partition(left, leftPoint - 1)
    partition(rightPoint + 1, right)
  }

  // [left,right]
  function randint(left: number, right: number) {
    const amplitude = right - left
    return (((amplitude + 1) * Math.random()) | 0) + left
  }

  function swap(i: number, j: number): void {
    const tmp = arr[i]
    arr[i] = arr[j]
    arr[j] = tmp
  }
}

export {}

if (require.main === module) {
  const arr = [4, 3, 2, 5, 6, 7, 8, 3, 2, 4, 1]
  quickSort(arr, (a, b) => a - b)
  console.log(arr)
}
