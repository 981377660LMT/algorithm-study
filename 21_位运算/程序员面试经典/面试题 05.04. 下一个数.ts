// 给定一个正整数，找出与其二进制表达式中1的个数相同且大小最接近的那两个数
// num的范围在[1, 2147483647]之间；
// 如果找不到前一个或者后一个满足条件的正数，那么输出 -1。
function findClosedNumbers(num: number): number[] {}

console.log(findClosedNumbers(0b10))
// 输出：[4, 1] 或者（[0b100, 0b1]）
console.log(Number(0x7fffffff).toString(2).length)

// 我直接投降
