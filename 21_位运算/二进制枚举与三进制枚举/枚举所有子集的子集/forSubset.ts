// forSubset枚举某个状态的所有子集(子集的子集)
const state = 0b1101
for (let g1 = state; ~g1; g1 = g1 === 0 ? -1 : (g1 - 1) & state) {
  if (g1 === state || g1 === 0) continue
  const g2 = state ^ g1
  console.log(g1.toString(2), g2.toString(2))
}

export {}
