# 如果一个十进制数字不含任何前导零，且每一位上的数字不是 0 就是 1 ，那么该数字就是一个 十-二进制数
# 给你一个表示十进制整数的字符串 n ，返回和为 n 的 十-二进制数 的最少数目。
class Solution:
    def minPartitions(self, n: str) -> int:
        return int(max(n))


print(Solution().minPartitions("32"))
# 输出：3
# 解释：10 + 11 + 11 = 32
