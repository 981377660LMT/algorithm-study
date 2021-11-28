const minutes = 3

console.log(Array.from({ length: minutes, 0: 1 }, (_, k) => k))
console.log(Array.from({ length: minutes }, (_, k) => Array(k).fill(k)))
