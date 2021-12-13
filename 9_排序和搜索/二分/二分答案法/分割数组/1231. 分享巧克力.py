from typing import List

# 将数组`顺序`划分成k+1个子数组，每段和有一个最小值
# 求这个最小值的最大可能值
# 0 <= K < sweetness.length <= 10^4

# 二分答案
class Solution:
    def maximizeSweetness(self, sweetness: List[int], k: int) -> int:
        def check(mid: int) -> bool:
            """"每段至少mid 可以分的段数>=k+1"""
            cur, count = 0, 0
            for num in sweetness:
                cur += num
                if cur >= mid:
                    cur = 0
                    count += 1
            return count >= k + 1

        left, right = 1, sum(sweetness)
        while left <= right:
            mid = (left + right) >> 1
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right


print(Solution().maximizeSweetness(sweetness=[1, 2, 3, 4, 5, 6, 7, 8, 9], k=5))
# 输出：6
# 解释：你可以把巧克力分成 [1,2,3], [4,5], [6], [7], [8], [9]。
