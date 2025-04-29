# 给定两个字符串s1 和 s2，返回 使两个字符串相等所需删除字符的 ASCII 值的最小和 。


class Solution:
    def minimumDeleteSum(self, s1: str, s2: str) -> int:
        """
        动态规划 O(m·n)：
        dp[i][j] = s1[:i] 和 s2[:j] 最小删除 ASCII 和
        转移：
          dp[0][0] = 0
          dp[i][0] = dp[i-1][0] + ord(s1[i-1])
          dp[0][j] = dp[0][j-1] + ord(s2[j-1])
          若 s1[i-1] == s2[j-1]:
            dp[i][j] = dp[i-1][j-1]
          否则:
            dp[i][j] = min(
              dp[i-1][j] + ord(s1[i-1]),   # 删除 s1 的第 i 个字符
              dp[i][j-1] + ord(s2[j-1])    # 删除 s2 的第 j 个字符
            )
        最终返回 dp[m][n]。
        空间可优化到 O(min(m,n))，但二维也足够。
        """
        m, n = len(s1), len(s2)
        dp = [[0] * (n + 1) for _ in range(m + 1)]
        for i in range(1, m + 1):
            dp[i][0] = dp[i - 1][0] + ord(s1[i - 1])
        for j in range(1, n + 1):
            dp[0][j] = dp[0][j - 1] + ord(s2[j - 1])
        for i in range(1, m + 1):
            ci = ord(s1[i - 1])
            for j in range(1, n + 1):
                cj = ord(s2[j - 1])
                if ci == cj:
                    dp[i][j] = dp[i - 1][j - 1]
                else:
                    dp[i][j] = min(dp[i - 1][j] + ci, dp[i][j - 1] + cj)
        return dp[m][n]


if __name__ == "__main__":
    sol = Solution()
    tests = [
        ("sea", "eat", 231),  # 删除 's' 和 't'
        ("delete", "leet", 403),  # 删除 "dee" or other组合
        ("", "abc", ord("a") + ord("b") + ord("c")),
        ("abc", "", ord("a") + ord("b") + ord("c")),
        ("abc", "abc", 0),
    ]
    for s1, s2, expect in tests:
        res = sol.minimumDeleteSum(s1, s2)
        print(f"{s1!r}, {s2!r} -> {res} (expected {expect})")
