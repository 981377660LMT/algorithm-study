// 有 N 位扣友参加了微软与力扣举办了「以扣会友」线下活动。
// 主办方提供了 2*N 道题目，整型数组 questions 中每个数字对应了每道题目所涉及的知识点类型。
// 若每位扣友选择不同的一题，请返回被选的 N 道题目至少包含多少种知识点类型。
function halfQuestions(questions: number[]): number {
  const n = questions.length / 2
  const counter = new Map<number, number>()
  questions.forEach(num => counter.set(num, (counter.get(num) || 0) + 1))
  const frequency = [...counter.values()].sort((a, b) => b - a)

  let remain = n
  let res = 1
  for (const freq of frequency) {
    if (freq >= remain) break
    remain -= freq
    res++
  }

  return res
}
console.log(halfQuestions([1, 5, 1, 3, 4, 5, 2, 5, 3, 3, 8, 6]))
