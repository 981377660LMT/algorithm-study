# 1≤N≤105,
# N  只排成一排的奶牛，编号为 1 到 N
# 如果 FJ 安排超过 K 只编号连续的奶牛，那么这些奶牛就会罢工去开派对。
# 找到最合理的安排方案并计算 FJ 可以得到的最大效率。

# 每部分子序列连续长度最多为m

# FJ 有 5 只奶牛，效率分别为 1、2、3、4、5。
# FJ 希望选取的奶牛效率总和最大，但是他不能选取超过 2 只连续的奶牛。
# 因此可以选择第三只以外的其他奶牛，总的效率为 1 + 2 + 4 + 5 = 12。

from collections import deque
from itertools import accumulate


n, m = map(int, input().split())
nums = [0] * n
for i in range(n):
    nums[i] = int(input())
preSum = [0] + list(accumulate(nums))

'''
dp(i)表示前i头牛的合法选择中最大总价值的数值
用最后一头牛的状态切分状态
1. 最后一头牛不选 子集1的最大价值就是dp(i-1)
2. 最后一头牛选，以结尾连续的牛个数进行子集划分, 设S(i)是前缀和
    最大价值是
    max {
        dp(j-2) + S(i) - S(j-1),
        dp(j-3) + S(i) - S(j-2),
        ......
        dp(j-m-1) + S(i) - S(j-m)        
    }
    那么只需要求i位置前面长度为m的滑动窗口中，dp(j-2)-S(j-1)这个序列在窗口中的最大值即可
'''


def cal(j: int) -> int:
    return dp[j - 1] - preSum[j]


# 队头维护的是cal(i)的最大值
queue = deque([0])
dp = [0] * (n + 1)
for i in range(1, len(preSum)):
    while i - queue[0] > m:
        queue.popleft()
    if queue:
        dp[i] = max(dp[i - 1], cal(queue[0]) + preSum[i])
    while queue and cal(queue[-1]) <= cal(i):
        queue.pop()
    queue.append(i)

print(dp[-1])

