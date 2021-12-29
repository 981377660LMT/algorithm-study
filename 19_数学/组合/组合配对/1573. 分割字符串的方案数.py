MOD = int(1e9 + 7)
# 请你返回分割 s 的方案数，满足 s1，s2 和 s3 中字符 '1' 的数目相同。
class Solution:
    def numWays(self, s: str) -> int:
        n = len(s)
        ones = [i for i, bit in enumerate(s) if bit == '1']

        m = len(ones)
        if m % 3 != 0:
            return 0

        # 字符串 s 中的所有字符都为 0  可任意分割 comb(n-1,2)
        if m == 0:
            ways = (n - 1) * (n - 2) // 2
            return ways % MOD
        else:
            # 起始索引
            index1, index2 = m // 3, m // 3 * 2
            count1 = ones[index1] - ones[index1 - 1]
            count2 = ones[index2] - ones[index2 - 1]
            ways = count1 * count2
            return ways % MOD


print(Solution().numWays(s="10101"))
# 输出：4
# 解释：总共有 4 种方法将 s 分割成含有 '1' 数目相同的三个子字符串。
# "1|010|1"
# "1|01|01"
# "10|10|1"
# "10|1|01"
