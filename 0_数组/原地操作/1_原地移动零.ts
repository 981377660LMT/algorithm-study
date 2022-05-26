// 必须在原数组上操作，不能拷贝额外的数组。
// 尽量减少操作次数。

// 最优解:统一处理0
const moveZeroes = (nums: number[]) => {
  let nonzeroIndex = 0

  for (let i = 0; i < nums.length; i++) {
    if (nums[i] !== 0) {
      nums[nonzeroIndex] = nums[i]
      nonzeroIndex++
    }
  }

  for (let i = nonzeroIndex; i < nums.length; i++) {
    nums[i] = 0
  }

  return nums
}

console.log(moveZeroes([0, 0, 1]))
console.log(moveZeroes([0, 1, 0, 3, 12]))

export {}

// ps:
// 1.
// webpack的源码, 发现webpack封装的ArrayQueue类中,
// 实现的出队列方法dequeue在数组长度大于16时, 采用reverse+pop来代替shift.
// https://segmentfault.com/a/1190000039183308 benchmark
// 2.
// shift方法每次调用时, 都需要遍历一次数组, 将数组进行一次平移, 时间复杂度是O(n).
// 而pop方法每次调用时, 只需进行最后一个元素的处理, 时间复杂度是O(1).
// 3
// 数组的大小不同决定了它在堆里存放的位置
// 小的数组(我猜是放在年轻分代里)在执行移动元素的操作时,
// 其实在堆中只是移动了指针而已.当大小超过一定数值,
// 数组将会被放到一个用于存放大对象的大对象空间(一页一个对象),
// 而由于内存对齐的原因(大概是页对齐?)不能通过移动指针实现,
// 只能真实的在内存中移动数据,因此效率降低.
// https://stackoverflow.com/questions/27341352/why-does-a-a-nodejs-array-shift-push-loop-run-1000x-slower-above-array-length-87
