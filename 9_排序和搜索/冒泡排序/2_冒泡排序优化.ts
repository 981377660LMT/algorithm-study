// const bubbleSort = (arr: number[]) => {
//   if (arr.length <= 1) return
//   // i表示有多少个元素排好序了
//   for (let i = 0; i < arr.length - 1; ) {
//     let lastSwapIndex = 0
//     for (let j = 0; j < arr.length - 1 - i; j++) {
//       if (arr[j] > arr[j + 1]) {
//         ;[arr[j], arr[j + 1]] = [arr[j + 1], arr[j]]
//         lastSwapIndex = j + 1
//       }
//       i = arr.length - lastSwapIndex
//     }
//   }
// }
// const arr = [1, 4, 2, 5, 3, 6, 7]
// bubbleSort(arr)
// console.log(arr)
// export {}

// 这个方法有问题
