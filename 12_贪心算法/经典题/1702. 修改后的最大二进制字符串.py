# 00 -> 10   (0生1)
# 10 -> 01  (1丢到0后面)
# 请你返回执行上述操作任意次以后能得到的 最大二进制字符串


# 总结：
# We don't need touch the starting 1s, they are already good.

# For the rest part,
# we continually take operation 2,
# making the string like 00...00011...11

# Then we continually take operation 1,
# making the string like 11...11011...11.
class Solution:
    def maximumBinaryString(self, binary: str) -> str:
        if '0' not in binary:
            return binary
        n = len(binary)
        trailingOnes = binary.count('1', binary.index('0'))
        return '1' * (n - trailingOnes - 1) + '0' + '1' * trailingOnes


print(Solution().maximumBinaryString(binary="000110"))
# 输出："111011"
# 解释：一个可行的转换为：
# "000110" -> "000101"
# "000101" -> "100101"
# "100101" -> "110101"
# "110101" -> "110011"
# "110011" -> "111011"
