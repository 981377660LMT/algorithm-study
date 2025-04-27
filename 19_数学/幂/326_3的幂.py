# 3的幂/三的幂
# 打表

POW3 = set(3**i for i in range(50))


class Solution:
    def isPowerOfThree(self, n: int) -> bool:
        while n > 0 and n % 3 == 0:
            n //= 3
        return n == 1

    def isPowerOfThree2(self, n: int) -> bool:
        return n in POW3

    def isPowerOfThree3(self, n: int) -> bool:
        """
        O(1) time, O(1) space solution using the fact that
        the maximum power of 3 in signed 32-bit range is 3^19 = 1162261467.
        A positive n is a power of three iff it divides 1162261467 exactly.
        """
        # 3**19 = 1162261467
        return n > 0 and 1162261467 % n == 0
