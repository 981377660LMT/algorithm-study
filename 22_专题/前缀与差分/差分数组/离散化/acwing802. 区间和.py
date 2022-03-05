# 假定有一个`无限长`的数轴，数轴上每个坐标上的数都是 0。
# 现在，我们首先进行 n 次操作，每次操作将某一位置 x 上的数加 c。
# 接下来，进行 m 次询问，每个询问包含两个整数 l 和 r，你需要求出在区间 [l,r] 之间的所有数的和。
# −109≤x≤109,
# 1≤n,m≤105,
# −109≤l≤r≤109,
# −10000≤c≤10000

# 操作和询问的坐标都是10^9级别，会超空间

# 思路：把`要用到的`查询、更新的坐标排序判重，用下标表示原来的值
# 再用二分查找离散化后的值

from bisect import bisect_left, bisect_right
from itertools import accumulate


n, m = map(int, input().split())
adds = []
queries = []
allNums = set()

for _ in range(n):
    x, c = map(int, input().split())
    adds.append((x, c))
    allNums.add(x)
for _ in range(m):
    l, r = map(int, input().split())
    queries.append((l, r))
    allNums.add(l)
    allNums.add(r)


allNums = sorted(allNums)
mapping = {allNums[i]: i for i in range(len(allNums))}

nums = [0] * len(allNums)
for x, c in adds:
    nums[mapping[x]] += c

preSum = [0] + list(accumulate(nums))
for l, r in queries:
    lower = bisect_left(allNums, l)
    upper = bisect_right(allNums, r)
    print(preSum[upper] - preSum[lower])

