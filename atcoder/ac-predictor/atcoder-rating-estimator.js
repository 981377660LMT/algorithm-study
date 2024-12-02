/* eslint-disable no-lone-blocks */
/* eslint-disable no-restricted-properties */
/* eslint-disable prefer-exponentiation-operator */
/* eslint-disable camelcase */
/* eslint-disable @typescript-eslint/no-unused-vars */

// https://github.com/koba-e964/atcoder-rating-estimator/blob/gh-pages/atcoder_rating.js
// https://koba-e964.github.io/atcoder-rating-estimator/test-last.html

// atcoder 评级计算
//
// 用于估算 AtCoder 用户评级的脚本。
// 该脚本通过分析用户的历史表现（Performance History）来计算未正化（unpositivized）和正化（positivized）的评级。
//
// 函数功能概览
// - calcAlgRatingFromHistory(history)：根据历史表现计算未正数化评级。
// - calcAlgRatingFromLast(last, perf, ratedMatches)：根据上一次的评级、新的表现和参赛次数，计算新的未正数化评级。
// - calcRequiredPerformance(targetRating, history)：计算为了达到目标评级，需要在下一次比赛中取得的表现。
// - calcHeuristicRatingFromHistory(history)：计算AHC评级。
// - positivizeRating(rating)：将未正数化评级转换为正数化评级。
// - unpositivizeRating(rating)：将正数化评级转换回未正数化评级。

const _finf = _bigf(400)
function _bigf(n) {
  let pow1 = 1
  let pow2 = 1
  let numerator = 0
  let denominator = 0
  for (let i = 0; i < n; ++i) {
    pow1 *= 0.81
    pow2 *= 0.9
    numerator += pow1
    denominator += pow2
  }
  return Math.sqrt(numerator) / denominator
}

/**
 * 调整评级的作用，主要是根据参赛次数 n，对计算的评级进行修正，使得随着参赛次数的增加，评级计算更加稳定和准确。
 * @param {Number} [n] 参赛次数
 *
 * f(n) 的作用和意义
  调整评级的初始偏差：在用户初次参加比赛时，由于缺乏历史数据，评级可能不稳定。f(n) 函数提供了一个根据参赛次数 n 的调整值，帮助评级系统在初期更快地收敛到准确的评级。
  使评级逐渐稳定：随着参赛次数的增加，f(n) 的值会逐渐减少，表示对评级的调整逐渐减小。这反映了随着用户参加更多的比赛，评级计算应更多地依赖实际表现，而不是初始的偏差修正。
  防止过度波动：对于新用户或参赛次数较少的用户，f(n) 能够防止评级出现过大的波动，提供更合理的评级评估。
 */
function _f(n) {
  return ((_bigf(n) - _finf) / (_bigf(1) - _finf)) * 1200.0
}

/**
 * calculate unpositivized rating from performance history
 * @param {Number[]} [history] performance history with ascending order
 * @returns {Number} unpositivized rating
 */
function calcAlgRatingFromHistory(history) {
  const n = history.length
  let pow = 1
  let numerator = 0.0
  let denominator = 0.0
  for (let i = n - 1; i >= 0; i--) {
    pow *= 0.9
    numerator += Math.pow(2, history[i] / 800.0) * pow
    denominator += pow
  }
  return Math.log2(numerator / denominator) * 800.0 - _f(n)
}

/**
 * calculate unpositivized rating from last state
 * @param {Number} [last] last unpositivized rating
 * @param {Number} [perf] performance
 * @param {Number} [ratedMatches] count of participated rated contest
 * @returns {number} estimated unpositivized rating
 */
function calcAlgRatingFromLast(last, perf, ratedMatches) {
  if (ratedMatches === 0) return perf - 1200
  last += _f(ratedMatches)
  const weight = 9 - 9 * 0.9 ** ratedMatches
  const numerator = weight * 2 ** (last / 800.0) + 2 ** (perf / 800.0)
  const denominator = 1 + weight
  return Math.log2(numerator / denominator) * 800.0 - _f(ratedMatches + 1)
}

/**
 * calculate the performance required to reach a target rate
 * @param {Number} [targetRating] targeted unpositivized rating
 * @param {Number[]} [history] performance history with ascending order
 * @returns {number} performance
 */
function calcRequiredPerformance(targetRating, history) {
  let valid = 10000.0
  let invalid = -10000.0
  for (let i = 0; i < 100; ++i) {
    const mid = (invalid + valid) / 2
    const rating = Math.round(calcAlgRatingFromHistory(history.concat([mid])))
    if (targetRating <= rating) valid = mid
    else invalid = mid
  }
  return valid
}

/**
 * calculate unpositivized rating from performance history
 * @param {Number[]} [history] performance histories
 * @returns {Number} unpositivized rating
 */
function calcHeuristicRatingFromHistory(history) {
  const S = 724.4744301
  const R = 0.8271973364
  const qs = []
  for (const perf of history) {
    for (let i = 1; i <= 100; i++) {
      qs.push(perf - S * Math.log(i))
    }
  }
  qs.sort((a, b) => b - a)
  let num = 0.0
  let den = 0.0
  for (let i = 99; i >= 0; i--) {
    num = num * R + qs[i]
    den = den * R + 1.0
  }
  return num / den
}

/**
 * (-inf, inf) -> (0, inf)
 * @param {Number} [rating] unpositivized rating
 * @returns {number} positivized rating
 */
function positivizeRating(rating) {
  if (rating >= 400.0) {
    return rating
  }
  return 400.0 * Math.exp((rating - 400.0) / 400.0)
}

/**
 * (0, inf) -> (-inf, inf)
 * @param {Number} [rating] positivized rating
 * @returns {number} unpositivized rating
 */
function unpositivizeRating(rating) {
  if (rating >= 400.0) {
    return rating
  }
  return 400.0 + 400.0 * Math.log(rating / 400.0)
}

// 根据历史表现计算未正数化的评级
{
  const history = [2400, 2400, 2400, 2400, 2400]
  const unpositivizedRating = calcAlgRatingFromHistory(history)
  console.log('未正数化的评级为：', unpositivizedRating)
}

// 根据上一次评级和新的比赛表现计算新的未正数化评级
{
  const lastRating = 1700
  const newPerformance = 1900
  const ratedMatches = 5
  const newUnpositivizedRating = calcAlgRatingFromLast(lastRating, newPerformance, ratedMatches)
  console.log('新的未正数化评级为：', newUnpositivizedRating)
}

// 计算达到目标评级所需的表现
{
  const history = [1500, 1600, 1700, 1750, 1780]
  const targetRating = 1800
  const requiredPerformance = calcRequiredPerformance(targetRating, history)
  console.log('为了达到目标评级，需要取得的表现为：', requiredPerformance)
}

// 将未正数化评级转换为正数化评级
// 当未正数化评级可能为负数或低于 400 时，可以使用 positivizeRating 函数进行转换。
{
  const unpositivizedRating = 350
  const positivizedRating = positivizeRating(unpositivizedRating)
  console.log('正数化的评级为：', positivizedRating)
}

// 根据历史表现使用启发式方法计算评级
{
  const history = [1400, 1500, 1600, 1700, 1800]
  const heuristicRating = calcHeuristicRatingFromHistory(history)
  console.log('启发式计算的未正数化评级为：', heuristicRating)
}

// !综合示例：完整的评级更新流程
// 假设一个用户的初始历史表现如下：
{
  let history = [1500, 1520, 1550]

  let currentRating = calcAlgRatingFromHistory(history)
  console.log('当前的未正数化评级为：', currentRating)

  let positivizedRating = positivizeRating(currentRating)
  console.log('当前的正数化评级为：', positivizedRating)

  const newPerformance = 1600
  history.push(newPerformance)

  currentRating = calcAlgRatingFromHistory(history)
  console.log('更新后的未正数化评级为：', currentRating)

  positivizedRating = positivizeRating(currentRating)
  console.log('更新后的正数化评级为：', positivizedRating)
}
