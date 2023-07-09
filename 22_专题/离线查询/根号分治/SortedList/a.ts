const arrs: number[][] = []
const nums: number[][] = []
for (let i = 0; i < 1e6; i++) {
  arrs.push([i, i, i, i])
  nums.push([i, i, i, i])
}

console.time('reset')
for (let i = 0; i < 1e6; i++) {
  arrs[i] = []
}
console.timeEnd('reset')

console.time('length')
for (let i = 0; i < 1e6; i++) {
  nums[i].splice(0, nums[i].length)
}
console.timeEnd('length')
