// 合并两个有序数组
const mergeTwo = (arr1: number[], arr2: number[]) => {
  const res: number[] = []
  let i = 0
  let j = 0

  while (i < arr1.length && j < arr2.length) {
    if (arr1[i] < arr2[j]) {
      res.push(arr1[i])
      i++
    } else {
      res.push(arr2[j])
      j++
    }
  }

  // 连接剩余的元素，防止没有把两个数组遍历完整
  return [...res, ...arr1.slice(i), ...arr2.slice(j)]
}

console.log(mergeTwo([1, 2, 3, 8, 9], [7, 10]))

export {}
