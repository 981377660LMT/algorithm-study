// 书店老板知道一个秘密技巧，能抑制自己的情绪，可以让自己连续 minutes 分钟不生气，但却只能使用一次。
// 请你返回这一天营业下来，最多有多少客户能够感到满意。
// 1 <= X <= customers.length == grumpy.length <= 20000
// 不生气时的顾客总数 + 窗口X内挽留的因为生气被赶走的顾客数
function maxSatisfied(customers: number[], grumpy: number[], minutes: number): number {
  const n = customers.length
  let customersWhenNotGrumpy = 0
  let retain = 0
  let maxRetain = 0

  for (let i = 0; i < n; i++) {
    if (grumpy[i] === 0) customersWhenNotGrumpy += customers[i]
  }

  for (let i = 0; i < minutes; i++) {
    if (grumpy[i] === 1) retain += customers[i]
  }

  maxRetain = retain

  for (let i = minutes; i < n; i++) {
    retain += customers[i] * grumpy[i] - customers[i - minutes] * grumpy[i - minutes]
    maxRetain = Math.max(maxRetain, retain)
  }

  return customersWhenNotGrumpy + maxRetain
}

console.log(maxSatisfied([1, 0, 1, 2, 1, 1, 7, 5], [0, 1, 0, 1, 0, 1, 0, 1], 3))
