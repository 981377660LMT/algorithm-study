# 一条直线上等距离放置了 n 台路由器。
# 路由器自左向右从 1 到 n 编号。
# 第 i 台路由器到第 j 台路由器的距离为 | i - j | 。
# 每台路由器都有自己的信号强度，
# 第 i 台路由器的信号强度为 ai 。
# 所有与第 i 台路由器距离不超过 ai 的路由器可以收到第i台路由器的信号（注意，每台路由器都能收到自己的信号）。
# 问一共有多少台路由器可以收到至少k台不同路由器的信号。
n, k = list(map(int, input().split()))
nums = list(map(int, input().split()))
ranges = []
for i in range(n):
    # 注意不要超边界
    left = max(0, i - nums[i])
    right = min(n, i + nums[i])
    ranges.append((left, right))

diff = [0 for _ in range(n + 10)]
for i in range(n):
    diff[ranges[i][0]] += 1
    diff[ranges[i][1] + 1] -= 1

res = 0
curSum = 0

for i in range(n):
    curSum += diff[i]
    if curSum >= k:
        res += 1
print(res)
