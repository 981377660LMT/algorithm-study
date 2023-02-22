const a1 = new Uint16Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16])
const a2 = new Uint16Array(16)
a2.set(a1)
console.log(a1, a2)
for (let i = 0; i < a1.length; i++) {
  a2[i] = 0
}
console.log(a1, a2)
