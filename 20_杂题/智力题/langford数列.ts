// 41312432 或 23421314
// 有8个数，11223344
// 将其排列，要求结果满足：两个1之间有一个数，两个2之间有两个数，两个3之间有三个数，两个4之间有四个数。问这个结果是多少？
// 可以证明 n%4 = 1和n%4 =2一定没有解
function generateLangFordSequence(n: number): number[] {
  if (n % 4 === 1 || n % 4 === 2) return []
}

console.log(generateLangFordSequence(4))
