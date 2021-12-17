# 请你返回分割 s 的方案数，满足 s1，s2 和 s3 中字符 '1' 的数目相同。
class Solution:
    def numWays(self, s: str) -> int:
        ...


print(Solution().numWays(s="10101"))
# 输出：4
# 解释：总共有 4 种方法将 s 分割成含有 '1' 数目相同的三个子字符串。
# "1|010|1"
# "1|01|01"
# "10|10|1"
# "10|1|01"
