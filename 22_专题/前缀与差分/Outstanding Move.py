# 可以交换一次数组两个元素，求∑(i+1)*nums[i]最大值
from itertools import accumulate


class Solution:
    def solve(self, nums):
        n = len(nums)
        preSum = [0] + list(accumulate(nums))
        base = sum((i + 1) * nums[i] for i in range(n))

        res = base

        # 将位置i的元素换到j的增量
        for i in range(n):
            for j in range(n + 1):
                if i == j:
                    continue

                res = max(res, base + preSum[i] - preSum[j] - (i - j) * nums[i])

        return res


print(Solution().solve(nums=[5, 1, 2]))
