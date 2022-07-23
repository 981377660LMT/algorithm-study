// const PI = 3.1415927
// const n = parseInt(readline())
// const nums = readline().split(' ').map(Number)
// const res = []
// for (let index = 0; index < n; index++) {
//   const num = nums[index]
//   const tmp = 0.5 * (1 - Math.cos((2 * PI * index) / n))
//   res.push(Math.round(num * tmp))
// }
// print(...res)

let n = parseInt(readline())
const nums = readline().split(' ').map(Number)
const visited = new Set()
const points = []
for (let index = 0; index < 2 * n; index += 2) {
  const hash = `${nums[index]}-${nums[index + 1]}`
  if (!visited.has(hash)) {
    visited.add(hash)
    points.push([nums[index], nums[index + 1]])
  }
}

n = points.length
let res = 0
for (let i = 0; i < n; i++) {
  const [x0, y0] = points[i]
  for (let j = 0; j < n; j++) {
    if (j === i) continue
    const [x1, y1] = points[j]
    for (let k = 0; k < n; k++) {
      if (k === i || k === j) continue
      const [x2, y2] = points[k]
      const ij = [x1 - x0, y1 - y0]
      const ik = [x2 - x0, y2 - y0]
      const jk = [x2 - x1, y2 - y1]
      const a = ij[0] * ij[0] + ij[1] * ij[1]
      const b = ik[0] * ik[0] + ik[1] * ik[1]
      const c = jk[0] * jk[0] + jk[1] * jk[1]
      if (a + b === c && a !== 0 && b !== 0 && c !== 0) res++
    }
  }
}

print(res / 2)

// TODO
// 两堆石子 5 个 7 个
// 每次可以一堆拿走若干个 或者两堆拿走同样的石头
// 问先手如何必胜
