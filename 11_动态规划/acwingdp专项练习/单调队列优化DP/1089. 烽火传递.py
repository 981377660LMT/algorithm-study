# 在某两个城市之间有 n 座烽火台，每个烽火台发出信号都有一定的代价。
# 为了使情报准确传递，在连续 m 个烽火台中至少要有一个发出信号。
# 现在输入 n,m 和每个烽火台的代价，请计算在两城市之间准确传递情报所需花费的总代价最少为多少。

from collections import deque


n, m = map(int, input().split())
nums = list(map(int, input().split()))

queue = deque([0])
dp = [0] * n


for i in range(1, n):
    while queue and i - queue[0] > m:
        queue.popleft()
    # // dp[i]表示点燃第i个烽火台时的最小花费
    dp[i] = dp[queue[0]] + nums[i]
    while queue and dp[queue[-1]] >= dp[i]:
        queue.pop()
    queue.append(i)

# 最后m个里选
print(min(dp[i] for i in range(n - m, n)))

