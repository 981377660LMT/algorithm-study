// 所有运动员的成绩都不相同。
function findRelativeRanks(score: number[]): string[] {
  const res = score.slice()
  const medalByScore = new Map<number, number | string>()

  score
    .sort((a, b) => b - a)
    .map((value, index) => {
      if (index === 0) medalByScore.set(value, 'Gold Medal')
      else if (index === 1) medalByScore.set(value, 'Silver Medal')
      else if (index === 2) medalByScore.set(value, 'Bronze Medal')
      else if (index > 2) medalByScore.set(value, index + 1)
    })

  return res.map(score => medalByScore.get(score)!.toString())
}

console.log(findRelativeRanks([5, 4, 3, 2, 1]))
// 输出: ["Gold Medal", "Silver Medal", "Bronze Medal", "4", "5"]
