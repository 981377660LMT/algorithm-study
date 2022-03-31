# Every Sublist Containing Unique Element
from collections import Counter
from typing import List

# O(n^2)


class Solution:
    def solve(self, nums: List[int]) -> bool:
        def check(left: int, right: int) -> bool:
            """[left,right]这段是否存在unique元素"""
            if left >= right:
                return True

            counter = Counter(nums[left : right + 1])
            if min(counter.values()) > 1:
                return False

            splitIndex = left
            for index in range(left, right + 1):
                # 需要检查左边
                if counter[nums[index]] == 1:
                    if not check(splitIndex, index - 1):
                        return False
                    splitIndex = index + 1

            return check(splitIndex, right)

        return check(0, len(nums) - 1)


print(Solution().solve(nums=[0, 2, 4, 2, 0]))
