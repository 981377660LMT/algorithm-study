const order1 = new Uint32Array(1e6)
for (let i = 0; i < order1.length; i++) order1[i] = i
const order2 = Array(1e6)
for (let i = 0; i < order2.length; i++) order2[i] = i

const nums = Array.from({ length: 1e6 }, () => (Math.random() * 1e5) | 0)

console.time('merge trick')
order1.sort((a, b) => nums[a] - nums[b])
console.timeEnd('merge trick')

console.time('naive')

order2.sort((a, b) => nums[a] - nums[b])
console.timeEnd('naive')
