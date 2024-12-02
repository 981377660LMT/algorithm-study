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
   */
  getRankForPerformance(performance) {
    if (this.rankMemo.has(performance)) return this.rankMemo.get(performance)
    const res = this.ratings.reduce((val, APerf) => val + 1.0 / (1.0 + Math.pow(6.0, (performance - APerf) / 400.0)), 0.5)
    this.rankMemo.set(performance, res)
    return res
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
