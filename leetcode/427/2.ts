export {}

function maxRectangleArea(points: number[][]): number {
  const h = (x: number, y: number) => `${x},${y}`
  const n = points.length
  const set = new Set<string>()
  for (const [x, y] of points) {
    set.add(h(x, y))
  }

  let res = -1
  for (let i = 0; i < n; i++) {
    for (let j = i + 1; j < n; j++) {
      const [x1, y1] = points[i]
      const [x2, y2] = points[j]

      if (x1 !== x2 && y1 !== y2) {
        const xMin = Math.min(x1, x2)
        const xMax = Math.max(x1, x2)
        const yMin = Math.min(y1, y2)
        const yMax = Math.max(y1, y2)

        const h1 = h(xMin, yMin)
        const h2 = h(xMin, yMax)
        const h3 = h(xMax, yMin)
        const h4 = h(xMax, yMax)

        if (set.has(h1) && set.has(h2) && set.has(h3) && set.has(h4)) {
          const hasOther = (): boolean => {
            for (let k = 0; k < n; k++) {
              if (k === i || k === j) continue
              const [xx, yy] = points[k]
              if (xx >= xMin && xx <= xMax && yy >= yMin && yy <= yMax) {
                const p = h(xx, yy)
                if (p !== h1 && p !== h2 && p !== h3 && p !== h4) {
                  return true
                }
              }
            }
            return false
          }

          if (!hasOther()) {
            const curRes = (xMax - xMin) * (yMax - yMin)
            res = Math.max(res, curRes)
          }
        }
      }
    }
  }

  return res
}
