// 给定两个已排序的整数数组，请找出其共有元素。
/**
 * @param {number[]} arr1 - integers
 * @param {number[]} arr2 - integers
 * @returns {number[]}
 */
function intersect(arr1: number[], arr2: number[]): number[] {
  let i = 0
  let j = 0
  let res: number[] = []

  while (i < arr1.length && j < arr2.length) {
    if (arr1[i] === arr2[j]) {
      res.push(arr1[i])
      i++
      j++
    } else if (arr1[i] < arr2[j]) {
      i++
    } else {
      j++
    }
  }

  return res
}

console.log(intersect([1, 2, 2, 3, 4, 4], [2, 2, 4, 5, 5, 6, 2000]))

// [2,2,4]

export {}
