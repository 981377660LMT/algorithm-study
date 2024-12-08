export {}

function maxRectangleArea(_ps: number[][]): number {
  const xCoord = _ps.map(p => p[0])
  const yCoord = _ps.map(p => p[1])
  const h = (x: number, y: number) => `${x},${y}`
  const n = xCoord.length
  const points = new Array(n)
  for (let i = 0; i < n; i++) {
    points[i] = [xCoord[i], yCoord[i]]
  }

  const mapY: Map<number, number[]> = new Map()
  for (const [x, y] of points) {
    if (!mapY.has(y)) {
      mapY.set(y, [])
    }
    mapY.get(y)!.push(x)
  }

  const pairMap: Map<string, number[]> = new Map()
  for (const [y, xs] of mapY) {
    xs.sort((a, b) => a - b)
    for (let i = 0; i < xs.length - 1; i++) {
      const x1 = xs[i]
      const x2 = xs[i + 1]
      const key = h(x1, x2)
      if (!pairMap.has(key)) pairMap.set(key, [])
      pairMap.get(key)!.push(y)
    }
  }

  const mapX: Map<number, number[]> = new Map()
  for (const [x, y] of points) {
    if (!mapX.has(x)) mapX.set(x, [])
    mapX.get(x)!.push(y)
  }
  for (const [_, ys] of mapX) {
    ys.sort((a, b) => a - b)
  }

  function isVerticalAdjacent(x: number, y1: number, y2: number): boolean {
    const ys = mapX.get(x)!
    let l = 0
    let r = ys.length
    let i1 = -1
    let i2 = -1
    while (l <= r) {
      const mid = (l + r) >> 1
      if (ys[mid] === y1) {
        i1 = mid
        break
      } else if (ys[mid] < y1) {
        l = mid + 1
      } else {
        r = mid - 1
      }
    }
    if (i1 === -1) return false

    l = 0
    r = ys.length - 1
    while (l <= r) {
      const mid = (l + r) >> 1
      if (ys[mid] === y2) {
        i2 = mid
        break
      } else if (ys[mid] < y2) {
        l = mid + 1
      } else {
        r = mid - 1
      }
    }
    if (i2 === -1) return false

    return Math.abs(i1 - i2) === 1
  }

  let maxArea = -1

  for (const [key, ys] of pairMap) {
    if (ys.length < 2) continue
    ys.sort((a, b) => a - b)
    const [x1, x2] = key.split(',').map(Number)
    for (let i = 0; i < ys.length - 1; i++) {
      const y1 = ys[i]
      const y2 = ys[i + 1]
      const width = x2 - x1
      const height = y2 - y1
      if (width > 0 && height > 0) {
        if (isVerticalAdjacent(x1, y1, y2) && isVerticalAdjacent(x2, y1, y2)) {
          const area = width * height
          if (area > maxArea) {
            maxArea = area
          }
        }
      }
    }
  }

  return maxArea
}
