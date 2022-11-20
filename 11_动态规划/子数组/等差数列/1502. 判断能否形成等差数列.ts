// 如果可以重新排列数组形成等差数列，请返回 true ；否则，返回 false
function canMakeArithmeticProgression(arr: number[]): boolean {
  if (arr.length < 3) return true
  arr.sort((a, b) => a - b)
  return new Set([...zip(arr, arr.slice(1))].map(([a, b]) => a - b)).size === 1

  function* zip(arr1: number[], arr2: number[]) {
    for (let i = 0; i < Math.min(arr1.length, arr2.length); i++) {
      yield [arr1[i], arr2[i]]
    }
  }
}

console.log(canMakeArithmeticProgression([3, 5, 1]))
export {}
