/**
 * @param {number[]} arr1 - integers
 * @param {number[]} arr2 - integers
 * @returns {number[]}
 */
function intersect(arr1: number[], arr2: number[]): number[] {
  let set1 = new Set(arr1)
  let set2 = new Set(arr2)
  if (set1.size > set2.size) {
    ;[set1, set2] = [set2, set1]
  }
  return [...set1].filter(num => set2.has(num))
}

console.log(intersect([1, 2, 2, 3, 4, 4], [2, 2, 4, 5, 5, 6, 2000]))

export {}
