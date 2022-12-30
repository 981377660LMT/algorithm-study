# 1009. 十进制整数的反码

# 补码: ~x
# 反码: 11111...111 ^ x


class Solution:
    def bitwiseComplement(self, n: int) -> int:
        if n == 0:
            return 1
        mask = (1 << n.bit_length()) - 1  # bit_length() = 32 - clz32(n)
        return n ^ mask


print(Solution().bitwiseComplement(5))
