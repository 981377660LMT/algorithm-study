// sum(Counter(S)[i] for i in J)
/**
 * @param {string} jewels
 * @param {string} stones
 * @return {number}
 * 统计stones中有多少个字符出现在了jewels中
 */
var numJewelsInStones = function (jewels, stones) {
  const set = new Set(jewels)
  return stones.split('').reduce((pre, cur) => pre + set.has(cur), 0)
}
