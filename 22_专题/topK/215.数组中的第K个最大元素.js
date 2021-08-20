// O(N)：用快排变形最最最高效解决TopK问题
// O(NlogK)：大根堆（前K小）/小根堆（前K大）
// O(N): 对于数据范围有限的情况例如本题，可以直接计数排序O(N)高效解决～

// 优先队列很容易
// 这里用快排的性质
const findKLargest = (arr, k) => {
  const expectedIndex = arr.length - k
  const findCurIndex = (arr, l, r) => {
    const curIndex = partition(arr, l, r)
    if (curIndex === expectedIndex) {
      return arr[curIndex]
    } else if (curIndex > expectedIndex) {
      return findCurIndex(arr, 0, curIndex - 1)
    } else {
      return findCurIndex(arr, curIndex + 1, arr.length - 1)
    }
  }

  return findCurIndex(arr, 0, arr.length - 1)
}

// 三路快排parition后可以知道leftPoint是数组第几小的元素(因为midPoint把小于pivot的数都移到前面去了)
// 然后在左边找/在右边找
// 时间复杂度n+n/2+n/4_... ≈ O(n)
const partition = (arr, l, r) => {
  if (arr.length <= 1) return
  if (l >= r) return

  // 优化，随机取标定点，以解决近乎有序的列表
  const randIndex = randint(l, r)
  swap(arr, l, randIndex)

  const pivot = arr[l]
  let leftPoint = l
  let midPoint = l
  let rightPoint = r
  while (midPoint <= rightPoint) {
    if (arr[midPoint] < pivot) {
      swap(arr, leftPoint, midPoint)
      leftPoint++
      midPoint++
    } else if (arr[midPoint] > pivot) {
      swap(arr, midPoint, rightPoint)
      rightPoint--
    } else {
      midPoint++
    }
  }
  console.log(arr, leftPoint, midPoint, rightPoint, 'pivot', pivot)
  return leftPoint
}

/**
 * @description 生成[start,end]的随机整数
 */
const randint = (start, end) => {
  if (start > end) throw new Error('invalid interval')
  const amplitude = end - start
  return Math.floor((amplitude + 1) * Math.random()) + start
}

const swap = (arr, i, j) => ([[arr[i], arr[j]]] = [[arr[j], arr[i]]])

console.log(findKLargest([3, 2, 1, 5, 6, 4, 4, 5, 5, 6], 2))
