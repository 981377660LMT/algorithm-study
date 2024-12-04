/* eslint-disable no-restricted-properties */
/* eslint-disable prefer-exponentiation-operator */

// 这段代码实现了一套用于计算和提供用户表现（Performance）和评级（Rating）的系统，
// 涉及多个提供器（Provider）类，每个类都负责不同的功能。
// 这种系统常用于竞技编程平台或游戏中，用于计算玩家的排名、表现和评级。

/**
 * 根据 Elo 评级系统，计算用户的表现值（Performance），基于用户的排名和现有的评级列表.
 */
class EloPerformanceProvider {
  ranks
  ratings
  cap
  rankMemo = new Map()

  constructor(ranks, ratings, cap) {
    this.ranks = ranks
    this.ratings = ratings
    this.cap = cap
  }
  availableFor(userScreenName) {
    return this.ranks.has(userScreenName)
  }
  getPerformance(userScreenName) {
    if (!this.availableFor(userScreenName)) {
      throw new Error(`User ${userScreenName} not found`)
    }
    const rank = this.ranks.get(userScreenName)
    return this.getPerformanceForRank(rank)
  }
  getPerformances() {
    const performances = new Map()
    for (const userScreenName of this.ranks.keys()) {
      performances.set(userScreenName, this.getPerformance(userScreenName))
    }
    return performances
  }

  /**
   * 根据排名计算对应的 perf，使用二分查找。
   */
  getPerformanceForRank(rank) {
    let upper = 6144
    let lower = -2048
    while (upper - lower > 0.5) {
      const mid = (upper + lower) / 2
      if (rank > this.getRankForPerformance(mid)) upper = mid
      else lower = mid
    }
    return Math.min(this.cap, Math.round((upper + lower) / 2))
  }

  /**
   * 根据 Perf 计算对应的预测排名，使用 Elo 公式。
   *
   * Elo 评级系统是一种广泛应用于棋类、竞技游戏和其他竞争性活动中的评级系统。它通过比较玩家之间的比赛结果来动态调整他们的评级（Rating），从而反映出玩家的相对实力。Elo 系统的核心思想是：
     - 预期得分（Expected Score）：根据两位玩家的评级差异，计算出一位玩家在比赛中获胜的概率。
     - 实际得分（Actual Score）：比赛结束后，根据实际结果调整玩家的评级。
     - 评级更新（Rating Update）：根据预期得分和实际得分的差异，调整玩家的评级。
     主要步骤：

     1. 缓存检查：如果该 performance 已经计算过，则直接返回缓存中的 rank。
     2. 计算排名：
        - 对所有现有的评级（ratings）进行遍历。
        - 对每个现有评级 APerf，计算与给定 performance 的比较结果。
        - 使用公式 `1.0 / (1.0 + 6.0^((performance - APerf) / 400.0))` 累加预期得分。
     3. 缓存结果：将计算出的 rank 存入 rankMemo 以备后续快速访问。
     4. 返回排名：输出最终计算的 rank。
   */
  getRankForPerformance(performance) {
    if (this.rankMemo.has(performance)) return this.rankMemo.get(performance)
    const res = this.ratings.reduce((val, APerf) => val + 1.0 / (1.0 + Math.pow(6.0, (performance - APerf) / 400.0)), 0.5)
    this.rankMemo.set(performance, res)
    return res // 最终的 res 可以被视为该 performance 在所有用户中的预期排名分数
  }
}

class InterpolatePerformanceProvider {
  ranks
  maxRank
  rankToUsers
  baseProvider
  constructor(ranks, baseProvider) {
    this.ranks = ranks
    this.maxRank = getMaxRank(ranks)
    this.rankToUsers = getRankToUsers(ranks)
    this.baseProvider = baseProvider
  }
  availableFor(userScreenName) {
    return this.ranks.has(userScreenName)
  }
  getPerformance(userScreenName) {
    if (!this.availableFor(userScreenName)) {
      throw new Error(`User ${userScreenName} not found`)
    }
    if (this.performanceCache.has(userScreenName)) return this.performanceCache.get(userScreenName)
    let rank = this.ranks.get(userScreenName)
    while (rank <= this.maxRank) {
      const perf = this.getPerformanceIfAvailable(rank)
      if (perf !== null) {
        return perf
      }
      rank++
    }
    this.performanceCache.set(userScreenName, -Infinity)
    return -Infinity
  }
  performanceCache = new Map()
  getPerformances() {
    let currentPerformance = -Infinity
    const res = new Map()
    for (let rank = this.maxRank; rank >= 0; rank--) {
      const users = this.rankToUsers.get(rank)
      if (users === undefined) continue
      const perf = this.getPerformanceIfAvailable(rank)
      if (perf !== null) currentPerformance = perf
      for (const userScreenName of users) {
        res.set(userScreenName, currentPerformance)
      }
    }
    this.performanceCache = res
    return res
  }
  cacheForRank = new Map()
  getPerformanceIfAvailable(rank) {
    if (!this.rankToUsers.has(rank)) return null
    if (this.cacheForRank.has(rank)) return this.cacheForRank.get(rank)
    for (const userScreenName of this.rankToUsers.get(rank)) {
      if (!this.baseProvider.availableFor(userScreenName)) continue
      const perf = this.baseProvider.getPerformance(userScreenName)
      this.cacheForRank.set(rank, perf)
      return perf
    }
    return null
  }
}

class FixedPerformanceProvider {
  result
  constructor(result) {
    this.result = result
  }
  availableFor(userScreenName) {
    return this.result.has(userScreenName)
  }
  getPerformance(userScreenName) {
    if (!this.availableFor(userScreenName)) {
      throw new Error(`User ${userScreenName} not found`)
    }
    return this.result.get(userScreenName)
  }
  getPerformances() {
    return this.result
  }
}

class IncrementalAlgRatingProvider {
  unpositivizedRatingMap
  competitionsMap
  constructor(unpositivizedRatingMap, competitionsMap) {
    this.unpositivizedRatingMap = unpositivizedRatingMap
    this.competitionsMap = competitionsMap
  }
  availableFor(userScreenName) {
    return this.unpositivizedRatingMap.has(userScreenName)
  }
  async getRating(userScreenName, newPerformance) {
    if (!this.availableFor(userScreenName)) {
      throw new Error(`rating not available for ${userScreenName}`)
    }
    const rating = this.unpositivizedRatingMap.get(userScreenName)
    const competitions = this.competitionsMap.get(userScreenName)
    return Math.round(positivizeRating(calcAlgRatingFromLast(rating, newPerformance, competitions)))
  }
}

class ConstRatingProvider {
  ratings
  constructor(ratings) {
    this.ratings = ratings
  }
  availableFor(userScreenName) {
    return this.ratings.has(userScreenName)
  }
  async getRating(userScreenName, newPerformance) {
    if (!this.availableFor(userScreenName)) {
      throw new Error(`rating not available for ${userScreenName}`)
    }
    return this.ratings.get(userScreenName)
  }
}

class FromHistoryHeuristicRatingProvider {
  performancesProvider
  constructor(performancesProvider) {
    this.performancesProvider = performancesProvider
  }
  availableFor(userScreenName) {
    return true
  }
  async getRating(userScreenName, newPerformance) {
    const performances = await this.performancesProvider(userScreenName)
    performances.push(newPerformance)
    return Math.round(positivizeRating(calcHeuristicRatingFromHistory(performances)))
  }
}
