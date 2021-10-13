const arr = [1, , , 2]

// forEach
arr.forEach(i => console.log(i)) // 1 2

// map
console.log(arr.map(i => i * 2)) // [ 2, <2 empty items>, 4 ]

// for ... of
for (const i of arr) {
  console.log(i)
}
// 1
// undefined
// undefined
// 2

// spread
console.log([...arr]) // [ 1, undefined, undefined, 2 ]
