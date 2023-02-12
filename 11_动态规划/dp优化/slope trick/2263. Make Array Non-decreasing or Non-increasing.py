# 2263. Make Array Non-decreasing or Non-increasing
# !每次可以选择一个数，将其增加或减少1，问最少需要多少次操作，使得数组满足非递减或非递增
# n<=1e5


from typing import List
from SlopeTrick import SlopeTrick


class Solution:
    def convertArray(self, nums: List[int]) -> int:
        def helper(arr: List[int]) -> int:
            """将数组变为非递减的+1-1最少操作数
            dp[i][x] = min(dp[i-1][y] + abs(x-arr[i])) , y<=x
            """
            st = SlopeTrick()
            for num in arr:
                st.add_abs(num)
                st.clear_right()
            return st.query()[0]

        return min(helper(nums), helper(nums[::-1]))


print(Solution().convertArray(nums=[3, 2, 4, 5, 0]))
