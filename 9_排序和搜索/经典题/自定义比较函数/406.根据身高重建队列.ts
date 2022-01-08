/**
 * @param {number[][]} people
 * @return {number[][]}
 * 每个 people[i] = [hi, ki] 表示第 i 个人的身高为 hi ，前面 正好 有 ki 个身高大于或等于 hi 的人。
 * 请你重新构造并返回输入数组 people 所表示的队列
 * @summary 一般这种数对，还涉及排序的，根据第一个元素正向排序，根据第二个元素反向排序，或者根据第一个元素反向排序，根据第二个元素正向排序，往往能够简化解题过程。
 * 如果两个维度一起考虑一定会顾此失彼。
 */
const reconstructQueue = function (people: number[][]): number[][] {
  const res: number[][] = []
  people.sort((a, b) => -(a[0] - b[0]) || a[1] - b[1])
  console.log(people)
  for (const p of people) {
    res.splice(p[1], 0, p)
  }
  return res
}

console.log(
  reconstructQueue([
    [7, 0],
    [4, 4],
    [7, 1],
    [5, 0],
    [6, 1],
    [5, 2],
  ])
)
// [[5,0],[7,0],[5,2],[6,1],[4,4],[7,1]]
export default 1
