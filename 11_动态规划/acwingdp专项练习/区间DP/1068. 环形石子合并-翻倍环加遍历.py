# 将 n 堆石子绕圆形操场排放，现要将石子有序地合并成一堆。

# 规定每次只能选相邻的两堆合并成新的一堆，并将新的一堆的石子数记做该次合并的得分。

# 请编写一个程序，读入堆数 n 及每堆的石子数，并进行如下计算：

# 选择一种合并石子的方案，使得做 n−1 次合并得分总和最大。
# 选择一种合并石子的方案，使得做 n−1 次合并得分总和最小。

# 1≤n≤200

# 本题要求每轮合并的石子 必须是相邻的 两堆石子，因此不能采用 Huffman Tree 的模型
# 这类限制只能合并相邻两堆石子的模型，用到的是经典的 区间DP 模型

from functools import lru_cache
from itertools import accumulate
from typing import Tuple


n = int(input())
nums = list(map(int, input().split()))

nums = nums * 2
preSum = [0] + list(accumulate(nums))


@lru_cache(None)
def dfs(left: int, right: int) -> Tuple[int, int]:
    """[left:right+1]这段区间合并石子的最值"""
    if left >= right:
        return 0, 0
    max_, min_ = -int(1e20), int(1e20)
    for i in range(left, right):
        max_ = max(max_, dfs(left, i)[0] + dfs(i + 1, right)[0] + preSum[right + 1] - preSum[left])
        min_ = min(min_, dfs(left, i)[1] + dfs(i + 1, right)[1] + preSum[right + 1] - preSum[left])
    return max_, min_


# 把环形的展开成2*n长度的链！然后遍历每个i到i+n-1的区间
max_, min_ = -int(1e20), int(1e20)
for start in range(n):
    curMax, curMin = dfs(start, start + n - 1)
    max_ = max(max_, curMax)
    min_ = min(min_, curMin)
print(min_)
print(max_)

