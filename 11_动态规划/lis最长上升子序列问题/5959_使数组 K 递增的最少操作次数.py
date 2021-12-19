from typing import List
from bisect import bisect_right

INF = 0x7FFFFFFF

# `没想起来是LIS` 还以为是直接二分找  => 题目变通能力还差了点

# 非递减改成严格递增怎么做?
# 把每个数`减去其下标`，然后对所有正整数求最长非降子序列。
class Solution:
    def kIncreasing(self, arr: List[int], k: int) -> int:
        # 是最长上升子序列
        def helper(arr: List[int]) -> int:
            LIS = [arr[0]]

            # 可以取相等：使用bisectRight
            for num in arr[1:]:
                if num >= LIS[-1]:
                    LIS.append(num)
                else:
                    index = bisect_right(LIS, num)
                    LIS[index] = num
            return len(arr) - len(LIS)

        return sum(helper(arr[start::k]) for start in range(k))


print(Solution().kIncreasing(arr=[5, 4, 3, 2, 1], k=1))
print(Solution().kIncreasing(arr=[4, 1, 5, 2, 6, 2], k=2))
print(Solution().kIncreasing(arr=[4, 1, 5, 2, 6, 2], k=3))
print(
    Solution().kIncreasing(arr=[12, 6, 12, 6, 14, 2, 13, 17, 3, 8, 11, 7, 4, 11, 18, 8, 8, 3], k=1)
)
# 12
