# 令 dp[i][j] 代表 s[i:j+1] 的最短编码。要么就是 s[i:j+1]，要么表示为 k[子串]，要么是 dp[i][k]+dp[k+1][j]。
# 判断某个字符串 ss 是否能表示为 k[子串] 可以参考 lc459。


class Solution:
    def encode(self, s: str) -> str:
        """
        DP over substrings: dp[i][j] = shortest encoding for s[i:j+1].
        For each substring:
          1. Initialize dp[i][j] = raw substring.
          2. Try all splits k: dp[i][j] = min(dp[i][k] + dp[k+1][j]).
          3. If substr can be represented as multiple repeats of a smaller pattern:
             find pattern length p via (substr+substr).find(substr, 1),
             then encode as str(repeat) + '[' + dp[i][i+p-1] + ']'.
        Return dp[0][n-1].
        Time: O(n^3) substrings * O(n) split/search ⇒ O(n^4) worst-case, but n≤150 passes.
        """
        n = len(s)
        dp = [[""] * n for _ in range(n)]
        for i in range(n):
            dp[i][i] = s[i]
        for length in range(2, n + 1):
            for i in range(n - length + 1):
                j = i + length - 1
                substr = s[i : j + 1]
                best = substr
                for k in range(i, j):
                    left = dp[i][k]
                    right = dp[k + 1][j]
                    if len(left) + len(right) < len(best):
                        best = left + right
                idx = (substr + substr).find(substr, 1)  # Find pattern length
                if idx < length:
                    times = length // idx
                    pattern_enc = dp[i][i + idx - 1]
                    candidate = f"{times}[{pattern_enc}]"
                    if len(candidate) < len(best):
                        best = candidate
                dp[i][j] = best
        return dp[0][n - 1]


if __name__ == "__main__":
    sol = Solution()
    tests = [
        ("aaa", "aaa"),
        ("aaaaaaaaaa", "10[a]"),
        ("aabcaabcd", "2[aabc]d"),
        ("abbbabbbcabbbabbbc", "2[2[abbb]c]"),
    ]
    for s, expected in tests:
        res = sol.encode(s)
        print(f"{s!r} -> {res} (expected {expected})")
