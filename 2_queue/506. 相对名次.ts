// 所有运动员的成绩都不相同。
function findRelativeRanks(score: number[]): string[] {
  const res = score.slice()
  const scoreToMedal = new Map<number, number | string>()
  score.sort((a, b) => b - a)
  score.map((value, index) => {
    if (index === 0) scoreToMedal.set(value, 'Gold Medal')
    else if (index === 1) scoreToMedal.set(value, 'Silver Medal')
    else if (index === 2) scoreToMedal.set(value, 'Bronze Medal')
    else if (index > 2) scoreToMedal.set(value, index + 1)
  })
  return res.map(score => scoreToMedal.get(score)!.toString())
}

console.log(findRelativeRanks([5, 4, 3, 2, 1]))
// 输出: ["Gold Medal", "Silver Medal", "Bronze Medal", "4", "5"]
