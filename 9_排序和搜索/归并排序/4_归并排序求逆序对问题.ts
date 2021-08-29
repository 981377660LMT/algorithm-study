/**
 * @param {number[]} nums
 * @return {number}
 * @description 归并排序mergeTwo的时候用两个指针统计逆序对
 */
const reversePairs = function (nums: number[]): number {
  let ans = 0

  //
  /**
   *
   * @param arr1 左边的数组
   * @param arr2 右边的数组
   * @description 合并两个有序数组
   * 这两个数组在源数组中元素下标 为i<j的关系
   */
  const mergeTwo = (arr1: number[], arr2: number[]) => {
    const ll = arr1.length
    const rl = arr2.length
    const res: number[] = Array(ll + rl)
    let l = 0
    let r = 0
    for (let index = 0; index < ll + rl; index++) {
      // 超出的情况
      if (l >= ll) {
        res[index] = arr2[r]
      } else if (r >= rl) {
        res[index] = arr1[l++]
        // 后面的全都大
      } else if (arr1[l] > arr2[r]) {
        ans += ll - l
        res[index] = arr2[r++]
      } else {
        res[index] = arr1[l++]
      }
    }

    return res
  }

  // 分left/right 递归 合
  const mergeSort = (arr: number[]): number[] => {
    if (arr.length <= 1) return arr

    const mid = Math.floor(arr.length / 2)
    const left = arr.slice(0, mid)
    const right = arr.slice(mid)

    return mergeTwo(mergeSort(left), mergeSort(right))
  }

  mergeSort(nums)

  return ans
}

const foo = [7, 5, 6, 4]
console.log(reversePairs(foo))
console.log(foo)
export {}
