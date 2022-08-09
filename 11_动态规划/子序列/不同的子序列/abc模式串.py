# 寻找abc模式串


class Solution:
    def solve(self, string):
        dp = [1, 0, 0, 0]
        for char in string:
            index = ord(char) - ord("a") + 1
            # endswith[index]表示取前, endswith[index - 1]表示不取重复的前
            dp[index] += dp[index] + dp[index - 1]
        return dp[-1]


print(Solution().solve("aabc"))
# We can make 2 "abc" and 1 "aabc"
