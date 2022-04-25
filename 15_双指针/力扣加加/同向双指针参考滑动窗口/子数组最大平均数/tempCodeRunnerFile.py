from itertools import accumulate
from sortedcontainers import SortedList

n, k = map(int, input().split())
nums = list(map(int, input().split()))


# 1 3 4
# preSum[i] - preSum[j] >= k * (i - j)

# 即 preSum[i]-k*i >= preSum[j] - k*j
preSum = [0] + list(accumulate(nums))
sortedList = SortedList()
res = 0
for i in range(n + 1):
    cur = preSum[i] - k * i
    pos = sortedList.bisect_right(cur)
    res += pos
    sortedList.add(cur)
print(res)