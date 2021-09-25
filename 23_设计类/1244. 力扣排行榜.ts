import { HashHeap } from '../8_heap/HashHeap'

// 排行榜是用Redis中zset有序集合实现的 而zset实现的一种则是跳表
class Leaderboard {
  private scoreRecord: Map<number, number>
  private scores: HashHeap<number> // 充当multiset

  constructor() {
    this.scoreRecord = new Map()
    this.scores = new HashHeap((a, b) => b - a)
  }

  addScore(playerId: number, score: number): void {
    if (!this.scoreRecord.has(playerId)) {
      this.scores.push(score)
      this.scoreRecord.set(playerId, score)
    } else {
      const preScore = this.scoreRecord.get(playerId)!
      this.scores.remove(preScore)
      this.scores.push(preScore + score)
      this.scoreRecord.set(playerId, preScore + score)
    }
  }

  // 返回前 K 名参赛者的 得分总和
  top(K: number): number {
    let res = 0
    while (this.scores.size && K--) {
      // res += this.scores.shift()!  // 不能shift
    }
    return res
  }

  // 将指定参赛者的成绩清零
  reset(playerId: number): void {
    if (this.scoreRecord.has(playerId)) return
    const score = this.scoreRecord.get(playerId)!
    this.scoreRecord.delete(playerId)
    this.scores.remove(score)
  }
}

export {}
// 1.一个map作为playId和对应score的映射，找到对应user的积分
// 2.另一个 multiset(排序好的集合，并且允许有相同的元素) 降序记录分数值

// 添加分数时：如果存在user 则删除multiset中user分数 重新添加新的分数
// topK计算分数直接遍历multiset前k大的分数
// reset先找到user分数 在multiset中删除分数
