# 1. 左端点从小到大排序
# 2. 遍历区间，判断能否判断放入现有组中

# 给定 N 个闭区间 [ai,bi]，
# 请你将这些区间分成若干组，
# 使得每组内部的区间两两之间（包括端点）没有交集，
# 并使得组数尽可能小。

# 输出最小组数。

#######################################################################
# 解答：
# 区间按照左端点排序
# 小根堆存放所有组的右端点值，堆顶存放最小的右端点值
# 如果当前区间左端点大于堆顶元素，说明可以加入堆顶元素所在组，右端点入堆
# 如果当前区间左端点小于等于堆顶元素，说明当前区间与堆里面的区间重叠

from heapq import heappop, heappush


n = int(input())

intervals = []
for _ in range(n):
    a, b = map(int, input().split())
    intervals.append((a, b))
intervals.sort()


pq = []
for start, end in intervals:
    if pq and start > pq[0]:
        heappop(pq)
    heappush(pq, end)

print(len(pq))

