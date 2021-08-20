/**
 * @param {number[]} ratings
 * @return {number}
 * 每个孩子至少分配到 1 个糖果。
   评分更高的孩子必须比他两侧的邻位孩子获得更多的糖果。  
   老师至少需要准备多少颗糖果呢?
   @summary 
   这道题目一定是要确定一边之后，再确定另一边，例如比较每一个孩子的左边，然后再比较右边，如果两边一起考虑一定会顾此失彼。
   从前向后遍历,只要右边评分比左边大，右边的孩子就多一个糖果
   再从后向前遍历,确定左孩子大于右孩子的情况
 */
const candy = function (ratings: number[]): number {
  const len = ratings.length
  const res: number[] = Array(len).fill(1)

  for (let i = 1; i < len; i++) {
    if (ratings[i] > ratings[i - 1]) {
      res[i] = res[i - 1] + 1
    }
  }

  for (let i = len - 1; i >= 1; i--) {
    if (ratings[i - 1] > ratings[i]) {
      res[i - 1] = Math.max(res[i - 1], res[i] + 1)
    }
  }

  return res.reduce((pre, cur) => pre + cur, 0)
}

console.log(candy([1, 0, 2]))
console.log(candy([1, 2, 2]))
// 输出：4
// 解释：你可以分别给这三个孩子分发 1、2、1 颗糖果。
//      第三个孩子只得到 1 颗糖果，这已满足上述两个条件。
export default 1
