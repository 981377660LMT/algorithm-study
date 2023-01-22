# 4的幂/四的幂


POW4 = set(4**i for i in range(32))


class Solution:
    def isPowerOfFour(self, n: int) -> bool:
        # 是二的倍数且减1是3的倍数
        return n > 0 and (n & (n - 1)) == 0 and (n - 1) % 3 == 0

    def isPowerOfFour2(self, n: int) -> bool:
        return n in POW4
