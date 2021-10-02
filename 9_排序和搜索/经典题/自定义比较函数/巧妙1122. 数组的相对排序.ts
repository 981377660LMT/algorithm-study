/**
 * @param {number[]} arr1
 * @param {number[]} arr2 
 * arr2 中的元素各不相同
   arr2 中的每个元素都出现在 arr1 中
 * @return {number[]}
 * 对 arr1 中的元素进行排序，使 arr1 中项的相对顺序和 arr2 中的相对顺序相同。
 * 未在 arr2 中出现过的元素需要按照升序放在 arr1 的末尾。
 * @summary
 * 1.按照arr2的顺序 那么就要比较函数为arr2中元素原来的顺序 即index
 * 2.不在arr2中统一加上arr2.length 使权重足够大以排到后面；这些元素之间又加上自身保持原来的大小
 */
const relativeSortArray = function (arr1: number[], arr2: number[]): number[] {
  const n = arr2.length
  const lookup = new Map<number, number>()
  arr2.forEach((num, index) => lookup.set(num, index))
  return arr1.sort((a, b) => {
    a = lookup.has(a) ? lookup.get(a)! : n + a
    b = lookup.has(b) ? lookup.get(b)! : n + b
    return a - b
  })
}

console.log(relativeSortArray([2, 3, 1, 3, 2, 4, 6, 7, 9, 2, 19], [2, 1, 4, 3, 9, 6]))
// 输出：[2,2,2,1,4,3,3,9,6,7,19]
