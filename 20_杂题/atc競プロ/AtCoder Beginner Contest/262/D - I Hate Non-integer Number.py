# n<=100 ai<=1e9
# !求平均数为整数的非空子集个数

# 确定选择的元素个数之后dp 因为限定了元素个数才知道mod是否为0(平均数为整数)
# !dp(index,remain,mod)
# !O(n^4)

# !启示:dp比记忆化dfs快两倍左右 时间卡的紧最好dp
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def main() -> None:
    # !TLE
    # def cal(select: int) -> int:
    #     def dfs(index: int, remain: int, mod: int) -> int:
    #         if remain < 0:
    #             return 0
    #         if index == n or remain == 0:
    #             return 1 if (remain == 0 and mod == 0) else 0
    #         hash_ = index * (n + 1) * (n + 1) + remain * (n + 1) + mod
    #         if memo[hash_] != -1:
    #             return memo[hash_]

    #         res = dfs(index + 1, remain, mod)
    #         if remain:
    #             res += dfs(index + 1, remain - 1, (mod + nums[index]) % select)
    #         res %= MOD
    #         memo[hash_] = res
    #         return res

    #     memo = [-1] * (n + 1) * (n + 1) * (n + 1)
    #     return dfs(0, select, 0)

    def cal(select: int) -> int:
        dp = [[0] * (select + 5) for _ in range(n + 5)]
        dp[0][0] = 1
        dp[1][nums[0] % select] = 1
        for i in range(1, n):
            ndp = [list(row) for row in dp]  # 不选
            for preS in range(i + 1):
                for preM in range(select):
                    ndp[preS + 1][(preM + nums[i]) % select] += dp[preS][preM]
                    ndp[preS + 1][(preM + nums[i]) % select] %= MOD
            dp = ndp

        return dp[select][0]

    n = int(input())
    nums = list(map(int, input().split()))

    res = 0
    for select in range(1, n + 1):
        res += cal(select)
        res %= MOD
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
