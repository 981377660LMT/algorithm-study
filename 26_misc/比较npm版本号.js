// npm 版本号

/**
 * @param {number} path
 */
function foo(path) {}
// 如果数组包含undefined，会使用默认值
let arr2 = [null, 20]
let [c = 3, d = 4] = arr2
console.log(c) // 输出 3
console.log(d) // 输出 20
