/**
 * @param {number} n
 * @return {number}
 * 第 i 轮，每 i 个灯泡切换一次开关。 
 * 而第 n 轮，你只切换最后一个灯泡的开关。
   找出 n 轮后有多少个亮着的灯泡。
   @summary 
   如果一个数的因数的个数为奇数个，那么它最后一定是亮的，否则是关闭的
   平方数一定是亮着的
 */
const bulbSwitch = function (n: number): number {
  return parseInt(Math.sqrt(n).toString())
}

console.log(bulbSwitch(3))

export default 1
