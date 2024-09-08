# No.911 ラッキーソート(异或+数位dp)
# https://yukicoder.me/problems/no/911/editorial
# 给定一个所有元素都不同的数组.
# 你可以对数组进行任意一次操作：
# 从[L, R]中选择一个数，所有数异或上这个数.
# 问有多少种方案，使得数组变成严格递增.
#
# n<=2e5.
# 按位填，通过计算msb，发现四种情况：
# 1.必须1
# 2.必须0
# 3.0和1都不可以
# 4.0和1都可以

from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    N, L, R = map(int, input().split())
    A = list(map(int, input().split()))

    status = [-1] * 64  # 0: 0, 1: 1, -1: 0 or 1

    def calc_bitwise_status() -> bool:
        for a, b in zip(A, A[1:]):
            msb = (a ^ b).bit_length() - 1
            if (a >> msb) & 1 == 0:
                if status[msb] == 1:
                    return False
                status[msb] = 0
            else:
                if status[msb] == 0:
                    return False
                status[msb] = 1
        return True

    if not calc_bitwise_status():
        print(0)
        exit()

    # 根据每个位选/不选，数位dp计算方案数
    @lru_cache(None)
    def dfs(v: int, bit: int) -> int:
        if v < 0:
            return 0
        if bit == 63:
            return 1
        s = status[bit]
        if s == 0:
            return dfs(v // 2, bit + 1)
        if s == 1:
            return dfs((v - 1) // 2, bit + 1)
        return dfs(v // 2, bit + 1) + dfs((v - 1) // 2, bit + 1)

    res1, res2 = dfs(R, 0), dfs(L - 1, 0)
    dfs.cache_clear()
    print(res1 - res2)
