// 思想：让数组越来越有序
// 对元素间距n/2的数组插入排序
// 对元素间距n/4的数组插入排序
// 对元素间距n/8的数组插入排序...
// 对元素间距1的数组插入排序...
// 克服插入排序一次只能交换相邻两个数的缺点
const shellSort = (nums: number[]) => {
  const len = nums.length
  let gap = 1
  while (gap < len / 3) {
    gap = gap * 3 + 1
  }
  //上面是设置动态增量算法 增量元素互质
  //下面是其实是插入排序 和 冒泡排序交换位置
  while (gap >= 1) {
    for (let i = 1; i < len; i++) {
      // 插入排序gap版
      for (let j = i; j >= gap && arr[j] < arr[j - gap]; j -= gap) {
        ;[arr[j], arr[j - gap]] = [arr[j - gap], arr[j]]
      }
    }
    gap = (gap - 1) / 3
  }
}

const arr = [4, 1, 2, 5, 3, 6, 7]
shellSort(arr)
console.log(arr)

export {}
// 希尔排序的基本思想是把数组按下标的一定增量分组，
// 对每组使用直接插入排序算法排序；随着增量逐渐减少，每组包含的元
// 素越来越多，当增量减至 1 时，整个数组恰被分成一组，算法便终止。
