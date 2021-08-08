/**
 * @param {number} jug1Capacity
 * @param {number} jug2Capacity
 * @param {number} targetCapacity
 * @return {boolean}
 */
const canMeasureWater = (
  jug1Capacity: number,
  jug2Capacity: number,
  targetCapacity: number
): boolean => {
  if (jug1Capacity + jug2Capacity < targetCapacity) return false
  const GCD = (a: number, b: number): number => (b === 0 ? a : GCD(b, a % b))
  return targetCapacity % GCD(jug1Capacity, jug2Capacity) === 0
}

console.log(canMeasureWater(3, 5, 4))

export {}
