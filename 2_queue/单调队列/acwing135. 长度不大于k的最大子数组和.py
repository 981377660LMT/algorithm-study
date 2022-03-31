# <!-- 输入一个长度为 n 的整数序列，从中找出一段长度不超过 m 的子数组
# ，使得子数组中所有数的和最大。 -->

# 1≤n,m≤300000
from collections import deque
from itertools import accumulate

'''
首先求序列的前缀和序列s, 将问题转换一下，以arr[i]结尾的长度不超过m的和最大
的连续子序列就是在s[i]前面的m个数中找最小的一个s[k]，s[i]-s[k]就是以arr[i]
结尾的长度不超过m的和最大的连续子序列的和，其实问题就转换成了单调队列求滑动
窗口极值问题

'''

n, m = map(int, input().split())
nums = list(map(int, input().split()))
preSum = [0] + list(accumulate(nums))

res = -int(1e20)
queue = deque([])


for i in range(len(preSum)):
    # 维护限制
    while queue and i - queue[0] > m:
        queue.popleft()
    if queue:
        res = max(res, preSum[i] - preSum[queue[0]])
    # 队首最小
    while queue and preSum[i] <= preSum[queue[-1]]:
        queue.pop()
    queue.append(i)

print(res)

