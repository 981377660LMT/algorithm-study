from math import ceil
from typing import List


MOD = int(1e9 + 7)
INF = int(1e20)

# !贪心+观察性质 拆木棒
# 0. 后面的数要尽可能大
# 1. 这个数要被分成几份? 二分求解或者ceil(cur/pre)
# 2. 把 x 分成 k 份，怎么让最小的那份最大？`divmod 平均分配`，最小的数为 div (即为 x//k)

# !divmod:平均分配 最大化最小值


class Solution:
    def minimumReplacement(self, nums: List[int]) -> int:
        n = len(nums)
        res = 0
        pre = INF
        for i in range(n - 1, -1, -1):
            cur = nums[i]
            if cur > pre:
                count = ceil(cur / pre)
                max_ = cur // count
                res += count - 1
                pre = max_
            else:
                pre = cur

        return res

    def minimumReplacement2(self, nums: List[int]) -> int:
        def cal(cur: int, pre: int) -> int:
            left, right = 2, int(1e9) + 5
            while left <= right:
                mid = (left + right) // 2
                if pre * mid >= cur:
                    right = mid - 1
                else:
                    left = mid + 1
            return left

        n = len(nums)
        res = 0
        pre = INF
        for i in range(n - 1, -1, -1):
            cur = nums[i]
            if cur > pre:
                count = cal(cur, pre)  # !需要拆成几个数 最大化拆出来的最小值 均分就是//n
                max_ = cur // count
                res += count - 1
                pre = max_
            else:
                pre = cur

        return res


print(Solution().minimumReplacement(nums=[3, 9, 3]))
print(Solution().minimumReplacement(nums=[1, 2, 3, 4, 5]))
