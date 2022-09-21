"""多重背包单调队列优化 时间复杂度O(n*v) 二进制的方法少一个log"""
# https://www.acwing.com/solution/content/6500/

# 有 N 种物品和一个容量是 V 的背包。
# 第 i 种物品`最多有 si 件`，每件体积是 vi，价值是 wi。
# 求解将哪些物品装入背包，可使物品体积总和不超过背包容量，且价值总和最大。
# 输出最大价值。

# 0<N≤1000
# 0<V≤20000
# 0<vi,wi,si≤20000

# 朴素的dp 为 dp[i][v] 表示前i种物品构成体积为v时的最大价值
# !单调队列优化的思路为 抛弃二维dp数组 只用一维数组dp[v] 表示体积为v时的最大价值
# 针对每一类物品 i ，我们都更新一下 dp[m] --> dp[0] 的值，最后 dp[m] 就是一个全局最优值
# dp[m] = max(dp[m], dp[m-v] + w, dp[m-2*v] + 2*w, dp[m-3*v] + 3*w, ...)

from collections import deque


n, cap = map(int, input().split())
dp = [0] * (cap + 1)
queue = deque()  # 单减的单调队列
for _ in range(n):
    cost, score, count = map(int, input().split())
    for j in range(cost):
        queue.clear()  # 注意每次新的循环都需要初始化队列
        remain = (cap - j) // cost  # 剩余的空间最多还能放几个当前物品
        for k in range(remain + 1):
            val = dp[k * cost + j] - k * score
            while queue and val >= queue[-1][1]:
                queue.pop()
            queue.append((k, val))
            while queue and queue[0][0] < k - count:  # 存放的个数不能超出物品数量，否则弹出
                queue.popleft()
            dp[k * cost + j] = queue[0][1] + k * score

print(dp[cap])
