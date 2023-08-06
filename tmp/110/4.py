from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个长度相等下标从 0 开始的整数数组 nums1 和 nums2 。每一秒，对于所有下标 0 <= i < nums1.length ，nums1[i] 的值都增加 nums2[i] 。操作 完成后 ，你可以进行如下操作：

# 选择任一满足 0 <= i < nums1.length 的下标 i ，并使 nums1[i] = 0 。
# 同时给你一个整数 x 。


# 请你返回使 nums1 中所有元素之和 小于等于 x 所需要的 最少 时间，如果无法实现，那么返回 -1 。


class Solution:
    def minimumTime(self, nums1: List[int], nums2: List[int], x: int) -> int:
        def check(mid: int) -> bool:
            """mid秒内能否使得nums1的和小于等于x."""
            sum_ = [a + mid * b for a, b in zip(nums1, nums2)]
            # 选择mid个数，减去mid*num2[a],(mid-1)*num2[a+1] ... 最大
            copy_ = nums2[:]
            copy_.sort(reverse=True)
            maxMinus = sum(copy_[i] * (mid - i) for i in range(mid))
            return sum(sum_) - maxMinus <= x

        left, right = 0, len(nums1)
        ok = False
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
                ok = True
            else:
                left = mid + 1
        return left if ok else -1


# nums1 = [1,2,3], nums2 = [1,2,3], x = 4
print(Solution().minimumTime(nums1=[1, 2, 3], nums2=[1, 2, 3], x=4))
