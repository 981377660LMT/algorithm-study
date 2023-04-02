# 0到x中二进制第k位为1的数的个数

from functools import lru_cache
from random import randint


def cal(upper: int, k: int) -> int:
    """[0,upper]中二进制第k(k>=0)位为1的数的个数.
    即满足 `num & (1 << k) > 0` 的数的个数
    """
    bit = upper.bit_length()
    if k >= bit:
        return 0

    @lru_cache(None)
    def dfs(pos: int, hasLeadingZero: bool, isLimit: bool) -> int:
        """当前在第pos位,hasLeadingZero表示有前导0，isLimit表示是否贴合上界"""
        if pos == -1:
            return not hasLeadingZero
        if pos == k:
            if isLimit and nums[pos] == 0:
                return 0
            return dfs(pos - 1, False, isLimit)
        res = 0
        up = nums[pos] if isLimit else 1
        for cur in range(up + 1):
            if hasLeadingZero and cur == 0:
                res += dfs(pos - 1, True, (isLimit and cur == up))
            else:
                res += dfs(pos - 1, False, (isLimit and cur == up))
        return res

    nums = list(map(int, bin(upper)[2:]))[::-1]
    res = dfs(len(nums) - 1, True, True)
    dfs.cache_clear()
    return res


if __name__ == "__main__":

    def bruteForce(upper: int, k: int) -> int:
        res = 0
        for i in range(upper + 1):
            res += (i & (1 << k)) > 0
        return res

    for i in range(100):
        upper, k = randint(0, int(1e5)), randint(0, 20)
        if cal(upper, k) != bruteForce(upper, k):
            print(upper, k)
            print(cal(upper, k), bruteForce(upper, k))
            break
