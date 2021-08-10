/**
 * @param {number[]} rains
 * @return {number[]}
 * @description 这道题没有使用到堆.事后诸葛亮这个技巧并不是堆特有的，实际上这就是一种普通的算法思想
 * 所有湖泊一开始都是空的。当第 n 个湖泊下雨的时候，如果第 n 个湖泊是空的，那么它就会装满水，否则这个湖泊会发生洪水。你的目标是避免任意一个湖泊发生洪水。
 * rains[i] > 0 表示第 i 天时，第 rains[i] 个湖泊会下雨。
 * rains[i] === 0 表示第 i 天没有湖泊会下雨，你可以选择 一个 湖泊并 抽干 这个湖泊的水。
 * 请返回一个数组 ans
 * 如果 rains[i] > 0 ，那么ans[i] == -1 。
   如果 rains[i] == 0 ，ans[i] 是你第 i 天选择抽干的湖泊。
   如果没办法阻止洪水，请返回一个 空的数组 。
   @summary 事后诸葛亮:晴天抽水存入集合，以后下雨了可以用来抵消
 */
const avoidFlood = function (rains: number[]): number[] {
  const res = Array<number>(rains.length).fill(1)

  // 记录晴天 注意这里用set是因为set插入时有序且便于删除
  const sunDay: Set<number> = new Set<number>()
  // 记录湖泊编号与下雨的天数
  const lakes = new Map<number, number>()

  for (let i = 0; i < rains.length; i++) {
    const lakeNum = rains[i]
    if (lakeNum > 0) {
      res[i] = -1

      if (lakes.has(lakeNum)) {
        let isFlood = true
        const rainDay = lakes.get(lakeNum)!
        // 这里需要寻找符合条件的晴天 因为可能是之前晴天然后两个雨天
        // 出现晴天的时候湖泊里面要有水才能抽
        // 找到第一个大于rainday的晴天，其实可以用二分
        for (const day of sunDay) {
          if (day > rainDay) {
            // 抽干水
            res[day] = rains[rainDay]
            sunDay.delete(day)
            lakes.delete(lakeNum)
            isFlood = false
            break
          }
        }
        if (isFlood) return []
      } else {
        lakes.set(lakeNum, i)
      }
    } else {
      sunDay.add(i)
    }
  }

  return res
}

console.log(avoidFlood([1, 0, 2, 0, 3, 0, 2, 0, 0, 0, 1, 2, 3]))
// console.log(avoidFlood([0, 1, 1]))
// 输出：[-1,-1,2,1,-1,-1]
