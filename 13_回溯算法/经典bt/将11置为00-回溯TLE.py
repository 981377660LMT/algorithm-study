from typing import List

# n ≤ 50


# TLE
class Solution:
    def solve(self, nums: List[int]) -> bool:
        """Assuming you play first and play optimally, return whether you can win the game."""

        def bt() -> bool:
            for i in range(len(nums) - 1):
                if nums[i] == nums[i + 1] == 1:
                    nums[i] = nums[i + 1] = 0
                    if not bt():
                        # 注意回溯 不要return就不管了
                        nums[i] = nums[i + 1] = 1
                        return True

                    nums[i] = nums[i + 1] = 1

            return False

        return bt()


print(Solution().solve([1, 1, 1, 1]))
print(Solution().solve([1, 1, 1, 1]))
