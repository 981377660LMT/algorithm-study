// 第i轮，不断找最小值，然后跟数组第i项交换
const selectSort = (arr: number[]) => {
  if (arr.length <= 1) return arr

  for (let i = 0; i < arr.length - 1; i++) {
    let minIndex = i
    for (let j = i; j < arr.length; j++) {
      if (arr[j] < arr[minIndex]) {
        minIndex = j
      }
    }
    ;[arr[minIndex], arr[i]] = [arr[i], arr[minIndex]]
  }
}

const arr = [4, 1, 2, 5, 3, 6, 7]
selectSort(arr)
console.log(arr)
export {}
