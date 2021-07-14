const m = 3

console.log(Array.from({ length: m, 0: 1 }, (_, k) => k))
console.log(Array.from({ length: m }, (_, k) => Array(k).fill(k)))
