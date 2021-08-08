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
  const queue: [number, number][] = [[0, 0]]
  const visited: Set<string> = new Set([`0#0`])

  const bfs = (): boolean => {
    while (queue.length) {
      const [x, y] = queue.shift()!
      if (x === targetCapacity || y === targetCapacity || x + y === targetCapacity) return true

      const aToB = Math.min(x, jug2Capacity - y)
      const bToA = Math.min(y, jug1Capacity - x)
      const next: [number, number][] = [
        [jug1Capacity, y],
        [x, jug2Capacity],
        [x, 0],
        [0, y],
        [x - aToB, y + aToB],
        [x + bToA, y - bToA],
      ]

      for (const [nextX, nextY] of next) {
        if (!visited.has(`${nextX}#${nextY}`)) {
          queue.push([nextX, nextY])
          visited.add(`${nextX}#${nextY}`)
        }
      }
    }

    return false
  }

  return bfs()
}

console.log(canMeasureWater(3, 5, 4))

export {}
