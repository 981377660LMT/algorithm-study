/**
 *
 * @param a
 * @param b
 * knows API 本质上还是在访问邻接矩阵
 */
declare function knows(a: number, b: number): boolean

/**
 * 
 * @param n  给你 n 个人的社交关系（你知道任意两个人之间是否认识），然后请你找出这些人中的「名人」。
 * 最多只有1个名人
 * @description
 * 所谓「名人」有两个条件：
   1、所有其他人都认识「名人」。
   2、「名人」不认识任何其他人。
   @summary
   名人节点的出度为 0，入度为 n - 1
 */
const findCelebrity = (n: number) => {
  // 假设0是名人
  let res = 0
  for (let other = 1; other < n; other++) {
    // res不可能是名人
    if (!knows(other, res) || knows(res, other)) {
      res = other
    }
  }

  // 现在的 res 是排除的最后结果，但不能保证一定是名人
  for (let other = 0; other < n; other++) {
    if (other === res) continue
    if (!knows(other, res) || knows(res, other)) return -1
  }

  return res
}
// 因为「名人」的定义保证了「名人」的唯一性，所以我们可以利用排除法，
// 先排除那些显然不是「名人」的人，从而避免 for 循环的嵌套，降低时间复杂度。
