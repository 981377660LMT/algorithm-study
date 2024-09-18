/**
 * @param {number[]} jobs  jobs[i] 是完成第 i 项工作要花费的时间。  1 <= k <= jobs.length <= 12
 * @param {number} k  k 位工人
 * @return {number}
 * 请你设计一套最佳的工作分配方案，使工人的 最大工作时间 得以 最小化
 * 返回分配方案中尽可能 最小 的 最大工作时间 。
 * @summary
 * 将一些数据块放在k个内存块。求内存块最小尺寸
 * 首次适应算法
 */
const minimumTimeRequired = function (jobs: number[], k: number): number {
  let res = Infinity
  // 剪枝1：先分配大的任务
  jobs.sort((a, b) => b - a)
  const workTime = Array<number>(k).fill(0)
  /**
   *
   * @param index job索引
   * @param workTime 各个工人的时间
   * @param maxTime 最大时间
   * @param used 分配给了多少个工人了
   * @returns
   */
  const bt = (index: number, workTime: number[], maxTime: number) => {
    if (maxTime >= res) return
    if (index === jobs.length) return (res = maxTime)

    // 剪枝2：相同的工作分配只使用第一次
    for (let i = 0; i < k; i++) {
      if (i >= 1 && workTime[i] === workTime[i - 1]) continue
      workTime[i] += jobs[index]
      bt(index + 1, workTime, Math.max(maxTime, workTime[i]))
      workTime[i] -= jobs[index]
    }
  }
  bt(0, workTime, 0)
  return res
}

console.log(minimumTimeRequired([1, 2, 4, 7, 8], 2))
// 输出：11
// 解释：按下述方式分配工作：
// 1 号工人：1、2、8（工作时间 = 1 + 2 + 8 = 11）
// 2 号工人：4、7（工作时间 = 4 + 7 = 11）
// 最大工作时间是 11 。

// 类似于473. 火柴拼正方形.ts
// 数据范围只有 12
// 爆搜（DFS）的复杂度是 O(k^n)
// 12^12 远超运算量 10^7

// 如何剪枝
// 1. 优先分配大的任务 jobs.sort((a,b)=>b-a)  4088 ms
// 优先分配工作量小的工作会使得工作量大的工作更有可能最后无法被分配
// 2. 相同的工作分配只使用第一次 if (i >= 1 && workTime[i] === workTime[i - 1]) continue  76 ms
// https://leetcode-cn.com/problems/find-minimum-time-to-finish-all-jobs/solution/gong-shui-san-xie-yi-ti-shuang-jie-jian-4epdd/
// 想要最大化剪枝效果，并且尽量让 k 份平均的话，我们应当调整我们对于「递归树」的搜索方向：将任务优先分配给「空闲工人」（带编号的方块代表工人）

export {}
