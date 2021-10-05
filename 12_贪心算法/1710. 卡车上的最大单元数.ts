/**
 * @param {number[][]} boxTypes
 * numberOfBoxesi 是类型 i 的箱子的数量。
   numberOfUnitsPerBoxi 是类型 i 每个箱子可以装载的单元数量。
 * @param {number} truckSize
 * @return {number}  返回卡车可以装载 单元 的 最大 总数。
 */
const maximumUnits = function (boxTypes: number[][], truckSize: number): number {
  boxTypes.sort((a, b) => b[1] - a[1])
  let res = 0
  for (const [num, size] of boxTypes) {
    const count = Math.min(num, truckSize)
    res += size * count
    truckSize -= count
    if (truckSize === 0) break
  }
  return res
}

console.log(
  maximumUnits(
    [
      [1, 3],
      [2, 2],
      [3, 1],
    ],
    4
  )
)
