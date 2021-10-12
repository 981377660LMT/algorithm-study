// 给定两个已排序的整数数组，请找出其共有元素。
/**
 * @param {number[]} arr1 - integers
 * @param {number[]} arr2 - integers
 * @returns {number[]}
 */
function intersect(arr1: number[], arr2: number[]): number[] {
  let p1 = 0
  let p2 = 0
  let result: number[] = []

  while (p1 < arr1.length && p2 < arr2.length) {
    if (arr1[p1] === arr2[p2]) {
      result.push(arr1[p1])
      p1++
      p2++
    } else {
      if (arr1[p1] < arr2[p2]) {
        p1++
      } else {
        p2++
      }
    }
  }

  return result
}

console.log(intersect([1, 2, 2, 3, 4, 4], [2, 2, 4, 5, 5, 6, 2000]))

// [2,2,4]

export {}
