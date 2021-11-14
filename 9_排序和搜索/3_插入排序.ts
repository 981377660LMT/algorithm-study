// 想象自己在打扑克整理手牌
// 每次抽出第i个元素 然后向前比较插入到合适的位置
// 如果待处理的数组本来就近乎有序，插入排序就是O(n)级别的复杂度
// 插入排序法升级:希尔排序法
const insertSort = (arr: number[]) => {
  if (arr.length <= 1) return

  for (let i = 1; i < arr.length; i++) {
    // for (let j = i; j >= 1; j--) {
    //   if (arr[j] < arr[j - 1]) {
    //     ;[arr[j], arr[j - 1]] = [arr[j - 1], arr[j]]
    //   } else {
    //     break
    //   }
    // }
    // 提前终止的机制
    for (let j = i; j >= 1 && arr[j - 1] > arr[j]; j--) {
      ;[arr[j], arr[j - 1]] = [arr[j - 1], arr[j]]
    }
  }
}

const arr = [4, 1, 2, 5, 3, 6, 7]
insertSort(arr)
console.log(arr)
export {}
