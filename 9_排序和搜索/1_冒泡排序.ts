// 相邻两个数不断比较，
// 互相交换冒泡到最后
// 第i论交换需要比较前0到arr.length-1-i个数
const bubbleSort = (arr: number[]) => {
  if (arr.length <= 1) return
  for (let i = 0; i < arr.length - 1; i++) {
    for (let j = 0; j < arr.length - 1 - i; j++) {
      if (arr[j] > arr[j + 1]) {
        ;[arr[j], arr[j + 1]] = [arr[j + 1], arr[j]]
      }
    }
  }
}
const arr = [1, 4, 2, 5, 3, 6, 7]
bubbleSort(arr)
console.log(arr)
export {}
