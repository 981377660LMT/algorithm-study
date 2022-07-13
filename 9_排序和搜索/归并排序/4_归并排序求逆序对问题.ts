/**
 * @param {number[]} nums
 * @return {number}
 * @description 归并排序mergeTwo的时候用两个指针统计逆序对
 */
function reversePairs(nums: number[]): number {
  let res = 0

  // 分left/right 递归 合
  const mergeSort = (arr: number[]): number[] => {
    if (arr.length <= 1) return arr

    const mid = Math.floor(arr.length / 2)
    const left = arr.slice(0, mid)
    const right = arr.slice(mid)

    return mergeTwo(mergeSort(left), mergeSort(right))
  }

  mergeSort(nums)

  /**
   *
   * @param arr1 左边的数组
   * @param arr2 右边的数组
   * @description 合并两个有序数组
   * 这两个数组在源数组中元素下标 为i<j的关系
   */
  function mergeTwo(arr1: number[], arr2: number[]) {
    const n1 = arr1.length
    const n2 = arr2.length
    const newArr: number[] = Array(n1 + n2).fill(0)
    let left = 0
    let right = 0
    for (let i = 0; i < n1 + n2; i++) {
      // 超出的情况
      if (left >= n1) {
        newArr[i] = arr2[right]
      } else if (right >= n2) {
        newArr[i] = arr1[left++]
        // 后面的全都大
      } else if (arr1[left] > arr2[right]) {
        res += n1 - left
        newArr[i] = arr2[right++]
      } else {
        newArr[i] = arr1[left++]
      }
    }

    return newArr
  }

  return res
}

const foo = [7, 5, 6, 4]
console.log(reversePairs(foo))
console.log(foo)
export {}
