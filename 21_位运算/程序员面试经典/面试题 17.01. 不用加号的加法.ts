/**
 * @param {number} a
 * @param {number} b
 * @return {number}
 */
const add = function (a: number, b: number): number {
  return a === 0 ? b : add((a & b) << 1, a ^ b)
}

console.log(add(1, 1))
export {}
