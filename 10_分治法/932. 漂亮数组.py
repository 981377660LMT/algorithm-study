#
# @param n 给定 N，返回任意漂亮数组 A（保证存在一个）。
# 数组 A 是整数 1, 2, ..., N 组成的排列
# 对于每个 i < j，都不存在 k 满足 i < k < j 使得 A[k] * 2 = A[i] + A[j]。
# @summary
# divide and conquer.
# put all even numbers on the left half, and all odd numbers on the right half, then recurse. O(n).
#
from typing import List


class Solution:
    def beautifulArray(self, n: int) -> List[int]:
        dp = [1]
        while len(dp) < n:
            ndp = []
            for num in dp:
                if num * 2 - 1 <= n:
                    ndp.append(num * 2 - 1)
            for num in dp:
                if num * 2 <= n:
                    ndp.append(num * 2)
            dp = ndp
        return dp


assert Solution().beautifulArray(4) == [1, 3, 2, 4]
