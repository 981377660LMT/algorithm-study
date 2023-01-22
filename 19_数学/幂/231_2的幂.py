# 2的幂/二的幂 : bit_count 等于1 (0b1000/0b0100/0b0010/0b0001...)

POW2 = set(1 << i for i in range(64))


class Solution:
    def isPowerOfTwo(self, i32: int) -> bool:
        return i32 > 0 and i32.bit_count() == 1
        return i32 > 0 and i32 & (i32 - 1) == 0

    def isPowerOfTwo2(self, i32: int) -> bool:
        return i32 in POW2


assert Solution().isPowerOfTwo(1)
