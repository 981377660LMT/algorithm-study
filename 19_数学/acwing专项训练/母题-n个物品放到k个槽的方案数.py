"""
n个数放到k个槽的方案数
"""
# !1.隔板法 从 n-1个空中选k-1个空 插隔板
# !如果不许有空槽 则为 C(n-1,k-1)
# !如果许有空槽 则为 C(n+k-1,k-1) 即 put(n,k)


from functools import lru_cache
from itertools import combinations


MOD = int(1e9 + 7)


@lru_cache(None)
def fac(n: int) -> int:
    """n的阶乘"""
    if n == 0:
        return 1
    return n * fac(n - 1) % MOD


@lru_cache(None)
def ifac(n: int) -> int:
    """n的阶乘的逆元"""
    return pow(fac(n), MOD - 2, MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac(n) * ifac(k)) % MOD * ifac(n - k)) % MOD


def put(n: int, k: int) -> int:
    """
    可以选取重复元素的组合数
    n个物品放入k个槽(槽可空)的方案数
    """
    return C(n + k - 1, k - 1)


# !2.dfs法 不太推荐 这样是n(logk)^2的dp
@lru_cache(None)
def dfs(k: int, n: int) -> int:
    """n个物品放到k个槽的放法 遍历每个槽计算 允许空槽"""
    if n < 0:
        return 0
    if k == 0:
        return 1 if n == 0 else 0
    res = 0
    for select in range(n + 1):  # 允许空槽
        res += dfs(k - 1, n - select)
        res %= MOD
    return res


if __name__ == "__main__":
    print(put(3, 2))
    print(dfs(2, 3))
    print(C(3, 2))
    print(len(list(combinations(range(3), 2))))
