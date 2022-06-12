// !数组元素分成k份 最小化求每份和的最大值
// 1 <= k <= jobs.length <= 12
function minimumTimeRequired(jobs: number[], k: number): number {
  let res = Infinity
  // 剪枝1：先分配大的任务
  jobs.sort((a, b) => b - a)
  bt(0, new Uint32Array(k), 0)
  return res

  function bt(index: number, workTime: Uint32Array, maxTime: number) {
    if (maxTime >= res) return
    if (index === jobs.length) {
      res = maxTime
      return
    }

    // 剪枝2：相同的工作分配只使用第一次
    for (let i = 0; i < k; i++) {
      if (i >= 1 && workTime[i] === workTime[i - 1]) continue
      workTime[i] += jobs[index]
      bt(index + 1, workTime, Math.max(maxTime, workTime[i]))
      workTime[i] -= jobs[index]
    }
  }
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
