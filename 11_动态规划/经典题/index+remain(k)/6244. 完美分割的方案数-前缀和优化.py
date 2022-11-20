# 给你一个字符串 s ，每个字符是数字 '1' 到 '9' ，再给你两个整数 k 和 minLength 。
# 如果对 s 的分割满足以下条件，那么我们认为它是一个 完美 分割：

# !1.s 被分成 k 段互不相交的子字符串。
# !2.每个子字符串长度都 至少 为 minLength 。
# !3.每个子字符串的第一个字符都是一个 质数 数字，最后一个字符都是一个 非质数 数字。
# 质数数字为 '2' ，'3' ，'5' 和 '7' ，剩下的都是非质数数字。

# 请你返回 s 的 完美 分割数目。由于答案可能很大，请返回答案对 1e9 + 7 取余 后的结果。
# 一个 子字符串 是字符串中一段连续字符串序列。

# !k,len(s),minLength<=1000
# 期望是O(n^2)的解法 所以需要优化掉dp范围转移的复杂度
# !注意到最内层的转移可以前缀和优化 => dp由一连串的index转移过来 所以考虑把index作为第二维度遍历

MOD = int(1e9 + 7)
ISPRIME = [False] * 10
for num in [2, 3, 5, 7]:
    ISPRIME[num] = True


class Solution:
    def beautifulPartitions(self, s: str, k: int, minLength: int) -> int:
        """dp[count][index]表示前index个字符分成count个子串的方案数"""
        nums = list(map(int, s))
        n = len(nums)
        dp = [[0] * (n + 1) for _ in range(k + 1)]
        dp[0][0] = 1

        for c in range(1, k + 1):
            dpSum = [0] * (n + 1)
            for i in range(1, n + 1):
                dpSum[i] = dpSum[i - 1] + (dp[c - 1][i - 1] if ISPRIME[nums[i - 1]] else 0)
                dpSum[i] %= MOD
            for i in range(1, n + 1):
                if not ISPRIME[nums[i - 1]]:
                    dp[c][i] = dpSum[max(0, i - minLength + 1)]  # i-length+1 个数的和

        return dp[k][n]


print(Solution().beautifulPartitions(s="23542185131", k=3, minLength=2))
