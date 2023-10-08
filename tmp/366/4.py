from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 和一个 正 整数 k 。

# 你可以对数组执行以下操作 任意次 ：

# 选择两个互不相同的下标 i 和 j ，同时 将 nums[i] 更新为 (nums[i] AND nums[j]) 且将 nums[j] 更新为 (nums[i] OR nums[j]) ，OR 表示按位 或 运算，AND 表示按位 与 运算。
# 你需要从最终的数组里选择 k 个元素，并计算它们的 平方 之和。

# 请你返回你可以得到的 最大 平方和。


# 由于答案可能会很大，将答案对 109 + 7 取余 后返回。


class Solution:
    def maxSum(self, nums: List[int], k: int) -> int:
        res = 0
        bitCounter = [0] * 32
        for num in nums:
            for i in range(32):
                if num & (1 << i):
                    bitCounter[i] += 1
        for _ in range(k):
            # 每个位上取走一个1
            cur = 0
            for i in range(32):
                if bitCounter[i] > 0:
                    cur |= 1 << i
                    bitCounter[i] -= 1
            if cur == 0:
                break
            res += cur * cur
            res %= MOD
        return res % MOD


# nums = [2,6,5,8], k = 2
print(Solution().maxSum([2, 6, 5, 8], 2))
# nums = [4,5,4,7], k = 3
print(Solution().maxSum([4, 5, 4, 7], 3))
