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
