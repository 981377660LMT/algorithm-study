// 分治算法

// 合并两个有序数组
const mergeTwo = (arr1: number[], arr2: number[]) => {
  const res: number[] = []

  // 如果任何一个数组为空，就退出循环
  while (arr1.length && arr2.length) {
    if (arr1[0] < arr2[0]) {
      res.push(arr1.shift()!)
    } else {
      res.push(arr2.shift()!)
    }
  }

  // 连接剩余的元素，防止没有把两个数组遍历完整
  return [...res, ...arr1, ...arr2]
}

// 分left/right 递归 合
const mergeSort = (arr: number[]): number[] => {
  if (arr.length <= 1) return arr

  const mid = Math.floor(arr.length / 2)
  const left = arr.slice(0, mid)
  const right = arr.slice(mid, arr.length)

  return mergeTwo(mergeSort(left), mergeSort(right))
}

const arr = [4, 1, 2, 5, 3, 6, 7]

console.log(mergeSort(arr))
export {}
