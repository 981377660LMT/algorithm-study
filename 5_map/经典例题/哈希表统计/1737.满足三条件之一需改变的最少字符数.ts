/**
 * @param {string} a
 * @param {string} b
 * @return {number}
 * 一步操作中，你可以将 a 或 b 中的 任一字符 改变为 任一小写字母 。
 *操作的最终目标是满足下列三个条件 之一 ：

a 中的 每个字母 在字母表中 严格小于 b 中的 每个字母 。
b 中的 每个字母 在字母表中 严格小于 a 中的 每个字母 。
a 和 b 都 由 同一个 字母组成。

 */
const minCharacters = function (a: string, b: string): number {
  const counterA = Array<number>(26).fill(0)
  const counterB = Array<number>(26).fill(0)
  console.log('z'.codePointAt(0)! - 97)
  for (const char of a) counterA[char.codePointAt(0)! - 97]++
  for (const char of b) counterB[char.codePointAt(0)! - 97]++
  const sum = (nums: number[]) => nums.reduce((pre, cur) => pre + cur, 0)
  const min = (nums: number[]) => Math.min.apply(null, nums)
  // console.log(counterA, counterB)
  // b中可取b到z 25种
  const way1 = Array.from(
    { length: 25 },
    (_, i) => sum(counterA.slice(i + 1)) + sum(counterB.slice(0, i + 1))
  )
  const way2 = Array.from(
    { length: 25 },
    (_, i) => sum(counterB.slice(i + 1)) + sum(counterA.slice(0, i + 1))
  )
  const way3 = Array.from({ length: 26 }, (_, i) => a.length + b.length - counterA[i] - counterB[i])
  console.log(way1, way2, way3)
  return min([min(way1), min(way2), min(way3)])
}

console.log(minCharacters('a', 'abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz'))

export default 1
