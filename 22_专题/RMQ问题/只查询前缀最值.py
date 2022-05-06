# 只查询前缀最值:O(n)处理即可
# 查询任意区间最值:O(nlogn)处理
# 给定n个任务需要花费的时间和产生的价值，最多只能花费m个时间，且只能完成两个任务，若不足2个任务，则返回0.
# 2 < n <= 1e6, 0 <= m <= 1e6

from bisect import bisect_right


n, m = map(int, input().split())
tasks = []
for _ in range(n):
    cost, value = map(int, input().split())
    tasks.append((cost, value))

tasks.sort()
preMax = [-int(1e20)]  # 动态更新

res = 0
for i, (cost, value) in enumerate(tasks):
    remainCost = m - cost
    pos = bisect_right(tasks, (remainCost, int(1e20))) - 1
    if pos >= 0:
        res = max(res, value + preMax[min(pos, len(preMax) - 1)])
    preMax.append(max(preMax[-1], value))

print(res)
