# 给定 N 个闭区间 [ai,bi] 以及一个线段区间 [s,t]，
# 请你选择尽量少的区间，将指定线段区间完全覆盖。
# 输出最少区间数，如果无法完全覆盖则输出 −1。

################################################
# 解法：
# 1.将所有区间按照左端点从小到大进行排序
# 2.从前往后枚举每个区间，在所有能覆盖start的区间中，
# 选择右端点的最大区间，然后`将start更新`成右端点的最大值

start, end = list(map(int, input().split()))

n = int(input())

intervals = []
for _ in range(n):
    a, b = map(int, input().split())
    intervals.append((a, b))
intervals.sort()

res = 0
index = 0
canCover = False
while index < n:
    curIndex = index
    curEnd = -int(1e20)

    while curIndex < n and intervals[curIndex][0] <= start:
        curEnd = max(curEnd, intervals[curIndex][1])
        curIndex += 1

    # 覆盖不了
    if curEnd < start:
        break

    # curIndex表示的区间是当前覆盖start的最右的区间
    res += 1

    # 已经覆盖玩
    if curEnd >= end:
        canCover = True
        break

    # 更新下一个start
    start = curEnd
    index = curIndex


print(res if canCover else -1)

