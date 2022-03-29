one = {str(i): 1 for i in range(1, 10)}
one.update({'*': 9, '0': 0})

two = {str(i): 1 for i in range(10, 27)}
two.update({'*' + str(i): 2 if i <= 6 else 1 for i in range(10)})
two.update({'1*': 9, '2*': 6, '**': 15})  # *表示1-9 所以 **表示(11-19,21-26)

MOD = 10 ** 9 + 7


class Solution:
    def numDecodings(self, s: str) -> int:
        if not s or s[0] == '0':
            return 0

        dp = [0 for _ in range(len(s) + 1)]
        dp[0] = 1
        dp[1] = one.get(s[0], 0)
        for i in range(2, len(s) + 1):
            dp[i] = (one.get(s[i - 1], 0) * dp[i - 1] + two.get(s[i - 2 : i], 0) * dp[i - 2]) % MOD
        return dp[-1]


print(Solution().numDecodings("2*"))
# 输入：s = "2*"
# 输出：15
# 解释：这一条编码消息可以表示 "21"、"22"、"23"、"24"、"25"、"26"、"27"、"28" 或 "29" 中的任意一条。
# "21"、"22"、"23"、"24"、"25" 和 "26" 由 2 种解码方法，但 "27"、"28" 和 "29" 仅有 1 种解码方法。
# 因此，"2*" 共有 (6 * 2) + (3 * 1) = 12 + 3 = 15 种解码方法。
