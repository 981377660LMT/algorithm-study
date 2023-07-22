// 2757. 生成循环数组的值
// https://leetcode.cn/problems/generate-circular-array-values/

// 给定你一个 循环 数组 arr 和一个整数 startIndex ，返回一个生成器对象 gen ，
// 它从 arr 中生成值。第一次调用 gen.next() 时，它应该生成 arr[startIndex] 。
// 每次调用 gen.next() 时，都会传入一个整数参数 jump（例如：gen.next(-3) ）。

// 如果 jump 是正数，则索引应该增加该值，但如果当前索引是最后一个索引，则应跳转到第一个索引。
// 如果 jump 是负数，则索引应减去该值的绝对值，但如果当前索引是第一个索引，则应跳转到最后一个索引。

function* cycleGenerator(arr: number[], startIndex: number): Generator<number, void, number> {
  const n = arr.length
  while (true) {
    const jump = yield arr[startIndex]
    startIndex = (startIndex + jump) % n
    if (startIndex < 0) startIndex += n
  }
}

/**
 *  const gen = cycleGenerator([1,2,3,4,5], 0);
 *  gen.next().value  // 1
 *  gen.next(1).value // 2
 *  gen.next(2).value // 4
 *  gen.next(6).value // 5
 */

export {}
