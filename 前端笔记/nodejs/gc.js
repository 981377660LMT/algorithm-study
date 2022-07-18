// let arr = []
// while (true) arr.push(1)
// invalid array length Allocation failed - JavaScript heap out of memory

// 没事
// let arr = []
// while (true) arr.push()

// !没事 说明Buffer对象的内存分配不是在V8的堆内存中，而是在Node的C++层面实现内存的申请的
// let arr = []
// while (true) arr.push(Buffer(1000))
