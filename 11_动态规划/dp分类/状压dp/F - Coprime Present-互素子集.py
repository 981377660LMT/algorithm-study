# 从 [a,b] (b-a ≤ 72,1<=a<=b<=1e18) 区间内找出不同的数构成一个子集，
# 子集中任意两数都互质，求这样不同的子集有多少个。
# !互素的子集数(包括空集)

# 1994. 好子集的数目-枚举因子
# !1. 注意gcd的差分性质 gcd(a,b) = gcd(a,b-a) <= b-a <=72
# !2. 注意72以内的质数很少,可以状压dp
# !3. 等价于72以内每个质数p只能在好子集中出现0或1次，对应着选或不选

from functools import lru_cache


P72 = [2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71]


def comprimePresent(A: int, B: int) -> int:
    @lru_cache(None)
    def dfs(index: int, state: int) -> int:
        if index == n:
            return 1
        res = dfs(index + 1, state)  # jump
        if state | masks[index] == state:
            res += dfs(index + 1, state ^ masks[index])
        return res

    n = B - A + 1
    # 每个数包含的质因子
    masks = [sum(1 << i for i, p in enumerate(P72) if num % p == 0) for num in range(A, B + 1)]
    res = dfs(0, (1 << len(P72)) - 1)
    dfs.cache_clear()
    return res


A, B = map(int, input().split())
print(comprimePresent(A, B))
