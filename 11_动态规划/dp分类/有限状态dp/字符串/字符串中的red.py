# https://www.nowcoder.com/discuss/1030672
# 长度n，小写字母组成的字符串中有多少个
# 至少含有2个red子串的字符串,结果取模1e9 + 7


# 状态机dp 可以画出状态转移矩阵(这样如果n很大还可以转化为矩阵快速幂)
# 思路是第二维来存，
# 0：无结尾，
# 1：以r结尾，
# 2：re结尾，
# 3：恰好含有一个red，
# 4：恰好含有一个red，且以r结尾
# 5：恰好含有一个red，且以re结尾
# 6：至少含有两个red

MOD = int(1e9 + 7)


def countRed(n: int) -> int:
    dp = [1, 0, 0, 0, 0, 0, 0]
    for _ in range(n):
        ndp = [0] * 7
        ndp[0] = (dp[0] * 25 + dp[1] * 24 + dp[2] * 24) % MOD
        ndp[1] = (dp[0] + dp[1] + dp[2]) % MOD
        ndp[2] = dp[1]
        ndp[3] = (dp[2] + dp[3] * 25 + dp[4] * 24 + dp[5] * 24) % MOD
        ndp[4] = (dp[3] + dp[4] + dp[5]) % MOD
        ndp[5] = dp[4]
        ndp[6] = (dp[5] + dp[6] * 26) % MOD
        dp = ndp
    return dp[-1]


print(countRed(7))
