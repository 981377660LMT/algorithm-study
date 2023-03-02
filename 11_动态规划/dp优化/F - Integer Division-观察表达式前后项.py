# https://www.cnblogs.com/Lanly/p/17094101.html
# F - Integer Division
# 分组/一维分割 => dp

# 给定一个数字字符串,现在可以将这个字符串分割成任意多个部分(n-1个分割点取或不取)
# 将分割的得分记为各个部分的乘积
# 求所有分割的得分的和模998244353的值
# n<=2e5


# 设 dp[i]表示前 i个数的的所有切割方案的代价和，转移就是枚举最后一个切割点位置。
# dp[0] = 1 dp[i] = ∑dp[j] * s[j+1:i] (0<=j<i)
# !如何优化成O(n)
# !考虑 dp[i]和dp[i+1]的转移式，发现两者非常相似
# dp[i] = ∑dp[j] * s[j+1:i] (0<=j<i)
# dp[i+1] = ∑dp[j] * s[j+1:i+1] (0<=j<i+1)
#         = ∑dp[j] * (10*s[j+1:i]+s[i+1]) (0<=j<i+1)
#         = 10*dp[i] + ∑dp[j] * s[i+1] (0<=j<i+1)

MOD = 998244353


def integerDivision(s: str) -> int:
    dp, preSum = 0, 1
    for i in range(len(s)):
        val = int(s[i])
        dp = (dp * 10 + preSum * val) % MOD
        preSum = (preSum + dp) % MOD
    return dp


if __name__ == "__main__":

    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    MOD = 998244353
    INF = int(4e18)
    n = int(input())
    x = input()
    print(integerDivision(x))
