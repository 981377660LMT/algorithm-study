from 康托展开 import calPerm


# 给定 n 和 k，返回第k个排列。
class Solution:
    def getPermutation(self, n: int, k: int) -> str:
        nums = list(range(1, n + 1))
        res = calPerm(nums, k - 1)
        return "".join(map(str, res))


assert Solution().getPermutation(3, 3) == "213"
