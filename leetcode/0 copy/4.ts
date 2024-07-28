export {}

const INF = 2e9 // !超过int32使用2e15
// # 给你一个二进制字符串 s。

// # 请你统计并返回其中 1 显著 的 子字符串 的数量。

// # 如果字符串中 1 的数量 大于或等于 0 的数量的 平方，则认为该字符串是一个 1 显著 的字符串

// !枚举长度
function numberOfSubstrings(s: string): number {}

console.time('a')
console.log(numberOfSubstrings('1'.repeat(4e4)))
console.timeEnd('a')
