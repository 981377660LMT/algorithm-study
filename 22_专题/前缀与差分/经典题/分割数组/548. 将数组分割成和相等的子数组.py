from typing import List
from itertools import accumulate

# 给定一个有 n 个整数的数组，你需要找到满足以下条件的三元组 (i, j, k) ：

# 0 < i, i + 1 < j, j + 1 < k < n - 1
# 子数组 (0, i - 1)，(i + 1, j - 1)，(j + 1, k - 1)，(k + 1, n - 1) 的和应该相等。
# 我们定义子数组 (L, R) 表示原数组从索引为L的元素开始至索引为R的元素。

# 1 <= n <= 2000。
# 1.每次固定中间的指针；
# 2.借助哈希 字典存储和加速查找；
# 3.前缀和，辅助计算子数组的和。


class Solution:
    def splitArray(self, nums: List[int]) -> bool:
        n = len(nums)
        preSum = [0] + list(accumulate(nums))

        for mid in range(3, n - 3):
            memo = set()
            for left in range(1, mid - 1):
                s1 = preSum[left]
                s2 = preSum[mid] - preSum[left + 1]
                if s1 == s2:
                    memo.add(s1)
            for right in range(mid + 2, n - 1):
                s3 = preSum[right] - preSum[mid + 1]
                s4 = preSum[n] - preSum[right + 1]
                if s3 == s4 and s3 in memo:
                    return True

        return False


print(Solution().splitArray([1, 2, 1, 2, 1, 2, 1]))
# 输出: True
# 解释:
# i = 1, j = 3, k = 5.
# sum(0, i - 1) = sum(0, 0) = 1
# sum(i + 1, j - 1) = sum(2, 2) = 1
# sum(j + 1, k - 1) = sum(4, 4) = 1
# sum(k + 1, n - 1) = sum(6, 6) = 1
