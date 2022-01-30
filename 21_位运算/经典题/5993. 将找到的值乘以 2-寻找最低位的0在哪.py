from typing import List

# 找出nums中最大的original*(2**c)的数


class Solution:
    def findFinalValue(self, nums: List[int], original: int) -> int:
        # mask 记录 nums 中含有哪些 original 的 2 幂次倍数。
        mask = 0
        for num in nums:
            if num % original == 0:
                k = num // original
                if k & (k - 1) == 0:
                    mask |= k

        # 找最低位的0，即取反后最低位的1
        mask = ~mask
        lowbit = mask & (-mask)
        return original * lowbit


print(Solution().findFinalValue(nums=[5, 3, 6, 1, 12], original=3))
print(Solution().findFinalValue(nums=[2, 7, 9], original=4))
