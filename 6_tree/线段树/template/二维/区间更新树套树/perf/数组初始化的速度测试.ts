// fill的速度测试-类型数组与普通数组
// !Benchmark: Array fill method vs for loop
// https://www.measurethat.net/Benchmarks/Show/9348/0/array-fill-method-vs-for-loop
// https://stackoverflow.com/questions/48836753/why-does-array-prototype-fill-have-such-a-large-performance-difference-compare
// 为什么Array.prototype.fill的性能差异如此之大？

let arrayTest = Array(1e7)

console.time('fill')
arrayTest.fill(1)
console.timeEnd('fill')

arrayTest = new Array(1e7)

console.time('for')
for (let i = 0; i < arrayTest.length; i++) {
  arrayTest[i] = 1
}

console.timeEnd('for')

// fill: 10.787s (其实与实现有关,有的环境快有的环境慢)
// for: 1.984s

// !对于普通数组,手动遍历赋值最快
// new Uint16Array(1).subarray()
// new Uint16Array(1).fill()
// new Uint16Array(1).set()
// new Uint16Array(1).copyWithin()
