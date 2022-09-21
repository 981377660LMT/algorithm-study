# 给定 N 个闭区间 [ai,bi]，请你在数轴上选择若干区间，
# 使得选中的区间之间互不相交（包括端点）。
# 输出可选取区间的最大数量。


################################################
# 1. 右端点从小到大排序
# 2. 遍历区间，如果已经包含点，pass，否则`选择当前区间右端点`

n = int(input())

intervals = []
for _ in range(n):
    a, b = map(int, input().split())
    intervals.append((a, b))


# 选择结束时间早的(罗志祥贪心问题)
intervals.sort(key=lambda x: x[1])

res = 0
preEnd = -int(1e20)
for start, end in intervals:
    if start > preEnd:
        res += 1
        preEnd = end

print(res)
