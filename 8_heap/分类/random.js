const maxLevel = 16
const maxRand = 2 ** maxLevel - 1
console.log(maxLevel - Math.log2(1 + Math.random() * maxRand))
console.log(~~4.8)
