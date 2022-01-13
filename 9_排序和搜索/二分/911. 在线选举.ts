import { bisectRight } from './7_二分搜索寻找最插右入位置'

class TopVotedCandidate {
  private times: number[]
  private winner: number[]

  /**
   *
   * @param persons
   * @param times
   * @description
   * 在选举中，第 i 张票是在时间为 times[i] 时投给 persons[i] 的。
   */
  constructor(persons: number[], times: number[]) {
    let votedPerson = -1
    let curMax = 0
    const counter = new Map<number, number>()
    const winner: number[] = []

    for (const p of persons) {
      counter.set(p, (counter.get(p) ?? 0) + 1)
      if (counter.get(p)! >= curMax) {
        curMax = counter.get(p)!
        votedPerson = p
      }
      winner.push(votedPerson)
    }

    this.times = times
    this.winner = winner
  }

  /**
   *
   * @param t
   * @returns 返回在 t 时刻主导选举的候选人的编号。
   */
  q(t: number): number {
    // console.log(this.winner)
    return this.winner[bisectRight(this.times, t) - 1]
  }
}

/**
 * Your TopVotedCandidate object will be instantiated and called as such:
 * var obj = new TopVotedCandidate(persons, times)
 * var param_1 = obj.q(t)
 */

const vote = new TopVotedCandidate([0, 1, 1, 0, 0, 1, 0], [0, 5, 10, 15, 20, 25, 30])

console.log(vote.q(3))
console.log(vote.q(12))
console.log(vote.q(25))
console.log(vote.q(15))
console.log(vote.q(24))
console.log(vote.q(8))
// [0,1,1,0,0,1]

export {}
