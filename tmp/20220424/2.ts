function countLatticePoints(circles: number[][]): number {
  let res = 0
  let [left, right, up, down] = [0, 210, 210, 0]
  for (const [x, y, r] of circles) {
    left = Math.min(left, x - r)
    right = Math.max(right, x + r)
    up = Math.max(up, y + r)
    down = Math.min(down, y - r)
  }

  for (let x = left; x <= right; x++) {
    for (let y = down; y <= up; y++) {
      for (const [cx, cy, cr] of circles) {
        if ((x - cx) ** 2 + Math.abs(y - cy) ** 2 <= cr ** 2) {
          res++
          break
        }
      }
    }
  }

  return res
}

// 卡python
