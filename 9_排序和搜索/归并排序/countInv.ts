import { SortedListFast } from '../../22_专题/离线查询/根号分治/SortedList/SortedListFast'

/**
 * 归并排序求逆序对.
 */
function countInvMergeSort(arr: ArrayLike<number>): number {
  if (arr.length < 2) return 0
  if (arr.length === 2) return arr[0] > arr[1] ? 1 : 0
  let res = 0
  let midCount = 0
  const upper: number[] = []
  const lower: number[] = []
  const mid = arr[Math.floor(arr.length / 2)]
  for (let i = 0; i < arr.length; i++) {
    const num = arr[i]
    if (num < mid) {
      lower.push(num)
      res += upper.length
      res += midCount
    } else if (num > mid) {
      upper.push(num)
    } else {
      midCount++
      res += upper.length
    }
  }
  res += countInvMergeSort(lower)
  res += countInvMergeSort(upper)
  return res
}

/**
 * SortedList求逆序对.
 */
function countInvSortedList(arr: ArrayLike<number>): number {
  const n = arr.length
  let res = 0
  const visited = new SortedListFast()
  for (let i = n - 1; ~i; i--) {
    const smaller = visited.bisectLeft(arr[i])
    res += smaller
    visited.add(arr[i])
  }
  return res
}

export { countInvMergeSort, countInvSortedList }

if (require.main === module) {
  const arr = [7, 5, 6, 4]
  console.log(countInvMergeSort(arr))
  console.log(countInvSortedList(arr))

  const bigArr = Array.from({ length: 1e6 }, () => Math.floor(Math.random() * 1e9))
  console.time('merge')
  console.log(countInvMergeSort(bigArr))
  console.timeEnd('merge')

  console.time('sortedList')
  console.log(countInvSortedList(bigArr))
  console.timeEnd('sortedList')
}
