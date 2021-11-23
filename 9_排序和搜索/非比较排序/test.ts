const arr1 = [1, 2, 3, 4, 5, 6, 7]
// for (const num of arr1) {
//   console.log(num)
//   arr1.push(99)
// }

// for (const num of arr1) {
//   console.log(num)
//   arr1.pop()
// }

// i会受到影响
for (let i = 0; i < arr1.length; i++) {
  console.log(i)
  arr1.pop()
}
