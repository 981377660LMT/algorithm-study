/**
 * @param {number} n
 * @param {number} m
 * @return {number}
 * 每次从这个圆圈里删除第m个数字（删除后从下一个数字开始计数）。求出这个圆圈里剩下的最后一个数字。
 * 3.数学解法
 * @link https://leetcode-cn.com/problems/yuan-quan-zhong-zui-hou-sheng-xia-de-shu-zi-lcof/solution/yuan-quan-zhong-zui-hou-sheng-xia-de-shu-j30k/
 * @description f(n,m)=(f(n-1,m)+m)%n
   f(n,m)指n个人，报第m个编号出列最终编号

   有0 1 2 3 4 5 6 7 8 9十个数字，假设m为3,最后结果可以先记成f(10,3)，即使我们不知道它是多少。
 */
const lastRemaining = (n: number, m: number): number => {
  if (n === 1) return 0
  return (lastRemaining(n - 1, m) + m) % n // 平移
}

console.log(lastRemaining(5, 3))
// console.log(lastRemaining(10, 17))

export default 1
