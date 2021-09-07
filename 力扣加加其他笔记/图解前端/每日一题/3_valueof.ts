// const a = {
//   count: 1,
//   valueOf() {
//     return this.count++
//   },
// }

// // @ts-ignore
// console.log(a == 1 && a == 2 && a == 3)

////////////////////////////////////////////////////////
let tmp = 1
Object.defineProperty(globalThis, 'a', {
  get() {
    return tmp++
  },
})
// @ts-ignore
console.log(a === 1 && a === 2 && a === 3)

export {}
