// 每次抽出第i个元素 然后向前比较插入到合适的位置
// 如果插入排序近乎有序，就是O(n)级别的复杂度
const insertSort = (arr: number[]) => {
  for (let i = 1; i < arr.length; i++) {
    for (let j = i; j > 0; j--) {
      if (arr[j - 1] > arr[j]) {
        ;[arr[j - 1], arr[j]] = [arr[j], arr[j - 1]]
      }
    }
  }
}

const arr = [4, 1, 2, 5, 3, 6, 7]
insertSort(arr)
console.log(arr)
export {}
