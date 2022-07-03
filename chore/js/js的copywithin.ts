// The copyWithin() method shallow copies part of an array to another location
// in the same array and returns it without modifying its length.

// The copyWithin works like C and C++'s memmove,
// and is a high-performance method to shift the data of an Array
// copyWithin 就像 C 和 C++ 的 memcpy 函数一样，且它是用来移动 Array 或者 TypedArray 数据的一个高性能的方法。
// 复制以及粘贴序列这两者是为一体的操作;即使复制和粘贴区域重叠，粘贴的序列也会有拷贝来的值。

// python切片/js的copywithin都是调用C的memmove

const arr = new Uint8Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
console.log(arr.copyWithin(0, -3, -1)) // 从(-3,-1)切片开始拷贝 从0开始覆盖
console.log(arr)

const arr2 = [1, 2]
console.log(arr2.reverse())
// start和end都是可以省略。
// start省略表示从0开始，end省略表示数组的长度值。

export {}

// class Queue<T> extends Array<T> {
//   override shift(): T | undefined {
//     const returnItem = this[0]
//     this.copyWithin(0, 1)  // 没什么用
//     this.pop()
//     return returnItem
//   }

//   override at(index: number): T | undefined {
//     if (index < 0) index += this.length
//     return this[index]
//   }
// }

// const deque = new Queue<number>()
// console.log(deque.push(1))
// console.log(deque.shift())
