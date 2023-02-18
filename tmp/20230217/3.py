from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 。

# 如果存在一些整数满足 0 <= index1 < index2 < ... < indexk < nums.length ，得到 nums[index1] | nums[index2] | ... | nums[indexk] = x ，那么我们说 x 是 可表达的 。换言之，如果一个整数能由 nums 的某个子序列的或运算得到，那么它就是可表达的。

# 请你返回 nums 不可表达的 最小非零整数 。


class Solution:
    def minImpossibleOR(self, nums: List[int]) -> int:
        S = set(nums)
        for i in range(33):
            if 1 << i not in S:
                return 1 << i


# [8388608,131072,128,2097152,65536,2048,438,1048576,8192,32,8,64,1024,2244,512,262144,4096,16384,4,256,2,4194304,2203,16,32768,410,524288,765,1]
# [4,32,16,8,8,75,1,2]
print(Solution().minImpossibleOR([4, 32, 16, 8, 8, 75, 1, 2]))
