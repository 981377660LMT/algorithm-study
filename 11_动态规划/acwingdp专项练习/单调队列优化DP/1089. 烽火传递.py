# 在某两个城市之间有 n 座烽火台，每个烽火台发出信号都有一定的代价。
# 为了使情报准确传递，在连续 m 个烽火台中至少要有一个发出信号。
# 现在输入 n,m 和每个烽火台的代价，请计算在两城市之间准确传递情报所需花费的总代价最少为多少。
# dp(i) 表示前i个烽火台进行选择，且最后一个烽火台点燃的所有方案中，最小开销的数值
# dp(i) = cost(i) + min(dp(i-1), dp(i-2), ... dp(i-m)) 其实dp(i)就依赖于前面长度为m
# 滑动窗口中的最小值
from collections import deque


n, m = map(int, input().split())
costs = list(map(int, input().split()))

queue = deque()
dp = [0] * n

for i in range(n):
    while queue and i - queue[0] > m:
        queue.popleft()
    # // dp[i]表示点燃第i个烽火台时的最小花费
    if i < m:
        dp[i] = costs[i]
    else:
        dp[i] = dp[queue[0]] + costs[i]
    while queue and dp[queue[-1]] >= dp[i]:
        queue.pop()
    queue.append(i)

# 最后m个里选
print(min(dp[-m:]))

