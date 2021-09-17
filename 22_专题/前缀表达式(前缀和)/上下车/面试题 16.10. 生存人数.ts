// 计算生存人数最多的年份。
// 如果有多个年份生存人数相同且均为最大值，输出其中最小的年份。
// 你可以假设所有人都出生于 1900 年至 2000 年（含 1900 和 2000 ）之间。
function maxAliveYear(birth: number[], death: number[]): number {
  const arr = Array<number>(102).fill(0)
  // eg:1900年出生的人导致1900年变化人数加1
  birth.forEach(up => arr[up - 1900]++)
  // eg:1900年死亡的人导致1901年变化人数减1
  death.forEach(down => arr[down + 1 - 1900]--)
  console.log(arr)
  let res = 0
  let max = -Infinity
  let sum = 0
  for (let i = 0; i < arr.length; i++) {
    sum += arr[i]
    if (sum > max) {
      res = 1900 + i
      max = sum
    }
  }
  return res
}

console.log(maxAliveYear([1900, 1901, 1950], [1948, 1951, 2000]))
