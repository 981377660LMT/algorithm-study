# 给定 N 个闭区间 [ai,bi]，请你在数轴上选择尽量少的点，使得每个区间内至少包含一个选出的点。
# 在数轴上选尽量少的点，使每个区间内至少包含一个选出的点

#############################################################
# 1. 右端点从小到大排序
# 2. 遍历区间，如果已经包含点，pass，否则`选择当前区间右端点`

INF = int(1e18)


n = int(input())

intervals = []
for _ in range(n):
    a, b = map(int, input().split())
    intervals.append((a, b))

intervals.sort(key=lambda x: x[1])

res = 0
preEnd = -INF
for start, end in intervals:
    if start > preEnd:
        res += 1
        preEnd = end
print(res)

# Hanging Banners
# !固定所有区间的最少点数
# !贪心思路：尽量在区间末尾钉钉子
