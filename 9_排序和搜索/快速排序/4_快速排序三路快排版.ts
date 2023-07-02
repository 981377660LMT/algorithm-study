/**
 *
 * @param arr
 * @returns
 * @description
 *（1）在数据集之中，选择一个元素作为"基准"（pivot）。
 *（2）所有小于"基准"的元素，都移到"基准"的左边；所有大于"基准"的元素，都移到"基准"的右边。
 *（3）对"基准"左边和右边的两个子集，不断重复第一步和第二步，直到所有子集只剩下一个元素为止。
 * @description 存在的问题:完全有序的数组 退化成 时间复杂度O(n^2) 空间复杂度O(n)
 * 解决:随机partion 生成[l,r]的随机值
 * ```js
 * const randIndex = randint(l, r)
  ;[[arr[l], arr[randIndex]]] = [[arr[randIndex], arr[l]]]
 * ```
 * @description 还是存在的问题 元素有很多是相同的 退化成 时间复杂度O(n^2) 空间复杂度O(n)
 * 解决:重新设计partition算法 双路快排  O(nlogn)
 * @description 还是存在的问题 双路快排对一样的元素处理了
 * 解决：三路快排
 */

/**
 * @description 开始选取pivot arr[l]
 * @description 小的放前面 等于的当中间 大的放后面
 * @description arr[l+1...p1-1]<v  arr[p1+1...p2-1]===v arr[p3...r]>v
 * @summary 三路快排的partition 多了一个指针将等于pivot的值放在中间
 * @summary 元素都相同的数组 O(n)
 */
const partition = (arr: number[], l: number, r: number) => {
  if (l < r) {
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

    partition(arr, l, leftPoint - 1)
    partition(arr, rightPoint + 1, r)

    return arr
  }
}

/**
 * @description 生成[start,end]的随机整数
 */
const randint = (start: number, end: number) => {
  if (start > end) throw new Error('invalid interval')
  const amplitude = end - start
  return Math.floor((amplitude + 1) * Math.random()) + start
}

function swap(arr: number[], i: number, j: number) {
  return ([[arr[i], arr[j]]] = [[arr[j], arr[i]]])
}

if (require.main === module) {
  const arr = [4, 3, 2, 5, 6, 7, 8, 3, 2, 4, 1]
  partition(arr, 0, arr.length - 1)
  console.log(arr)
}

export {}
