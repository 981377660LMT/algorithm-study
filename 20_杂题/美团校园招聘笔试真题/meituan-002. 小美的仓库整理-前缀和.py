from bisect import bisect_left
from collections import defaultdict
from itertools import accumulate


MOD = int(1e9 + 7)
INF = 0x3F3F3F3F
EPS = int(1e-8)
DIRS4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
DIRS8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]


n = int(input())
# 读取数组
nums = [int(v) for v in input().split()]
queries = [int(v) - 1 for v in input().split()]


preSum = [0] + list(accumulate(nums))
sortedList = list(range(-1, n + 1))
res = [0] * n
res[n - 1] = 0

for i in range(n - 1, 0, -1):
    id = bisect_left(sortedList, queries[i])
    left = sortedList[id - 1]
    right = sortedList[id + 1]
    sortedList.pop(id)

    cur = preSum[right] - preSum[left + 1]
    res[i - 1] = max(res[i], cur)


# 输出答案
for v in res:
    print(v)

