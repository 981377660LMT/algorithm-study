# 爬n阶楼梯 每次只能爬奇数级
# 有些楼梯断了不能爬
# 爬到最后一个位置有多少种方法
# !dp[i] = dp[i-1] + dp[i-3] + dp[i-5] + ...
# !令 ep[i] = dp[i] + dp[i-2] + dp[i-4] + ...
# !那么 dp[i] = ep[i-1] = dp[i-1] + ep[i-3]
# dp[i] ep[i-1] ep[i-2] 的关系可由矩阵快速幂 logS 求出
# 转移矩阵
# 1 0 1
# 1 0 1
# 0 1 0
# 初始向量 [1 0 0]

# !坏的楼梯处
# !矩阵快速幂中间要打断 每次都直接把dp[i]手动赋值为0 再继续算
# 总时间复杂度 O(nlogS)
from functools import lru_cache
import gc
import sys
import os
import numpy as np

gc.disable()

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = 998244353


NPArray = np.ndarray


def main() -> None:
    @lru_cache(None)
    def matqpow2(exp: int) -> NPArray:
        """矩阵快速幂np版"""
        if exp == 0:
            return np.eye(*trans.shape, dtype=np.uint64)
        if exp == 1:
            return trans.copy()

        if exp & 1:
            return (matqpow2(exp - 1) @ trans) % MOD
        half = matqpow2(exp // 2)
        return (half @ half) % MOD

    n, s = map(int, input().split())
    bad = list(map(int, input().split()))
    res = np.array([[1], [0], [0]], np.uint64)  # 3 x 1 答案矩阵
    trans = np.array([[1, 0, 1], [1, 0, 1], [0, 1, 0]], np.uint64)
    pre = 0
    for cur in bad:
        res = (matqpow2(cur - pre) @ res) % MOD
        res[0][0] = 0
        pre = cur

    res = (matqpow2(s - bad[-1]) @ res) % MOD
    print(int(res[0][0]))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
