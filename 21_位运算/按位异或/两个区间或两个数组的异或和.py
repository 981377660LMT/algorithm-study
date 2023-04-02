# 两个区间的异或和
# 给定0<=x<=n,0<=y<=m的所有(x,y)对,求异或和
# 0<=n,m<=1e9
# !按位统计

from functools import lru_cache
from random import randint
from typing import List

MOD = int(1e9 + 7)


def cal(upper: int, k: int) -> int:
    """[0,upper]中二进制第k(k>=0)位为1的数的个数.
    即满足 `num & (1 << k) > 0` 的数的个数
    """
    bit = upper.bit_length()
    if k >= bit:
        return 0

    @lru_cache(None)
    def dfs(pos: int, hasLeadingZero: bool, isLimit: bool) -> int:
        """当前在第pos位,hasLeadingZero表示有前导0,isLimit表示是否贴合上界"""
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


def xorSum1(n: int, m: int) -> int:
    """给定0<=x<=n,0<=y<=m的所有(x,y)对,求异或和模1e9+7"""
    res = 0
    for bit in range(33):
        count = cal(n, bit) * (m + 1 - cal(m, bit)) + cal(m, bit) * (n + 1 - cal(n, bit))
        count %= MOD
        res += (1 << bit) * count
        res %= MOD
    return res


########################################################
########################################################
def xorSum2(nums1: List[int], nums2: List[int]) -> int:
    """两个数组所有数对的异或和模1e9+7."""
    C1 = [0] * 33
    C2 = [0] * 33
    for num in nums1:
        for bit in range(33):
            if num & (1 << bit):
                C1[bit] += 1
    for num in nums2:
        for bit in range(33):
            if num & (1 << bit):
                C2[bit] += 1

    n1, n2 = len(nums1), len(nums2)
    res = 0
    for bit in range(33):
        count = C1[bit] * (n2 - C2[bit]) + C2[bit] * (n1 - C1[bit])
        count %= MOD
        res += (1 << bit) * count
        res %= MOD
    return res


if __name__ == "__main__":

    def bruteForce1(n: int, m: int) -> int:
        res = 0
        for i in range(n + 1):
            for j in range(m + 1):
                res += i ^ j
        return res

    print(xorSum1(int(1e7), int(1e7)))
    for _ in range(100):
        n, m = randint(1, 1000), randint(1, 1000)
        if xorSum1(n, m) != bruteForce1(n, m):
            print(n, m)
            print(xorSum1(n, m), bruteForce1(n, m))
            break

    def bruteForce2(nums1: List[int], nums2: List[int]) -> int:
        res = 0
        for num1 in nums1:
            for num2 in nums2:
                res += num1 ^ num2
        return res

    for _ in range(100):
        nums1 = [randint(1, 1000) for _ in range(10)]
        nums2 = [randint(1, 1000) for _ in range(10)]
        if xorSum2(nums1, nums2) != bruteForce2(nums1, nums2):
            print(nums1, nums2)
            print(xorSum2(nums1, nums2), bruteForce2(nums1, nums2))
            break
