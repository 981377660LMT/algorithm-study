"""
康托展开 - 无重复元素
求字典序第k小的排列/当前排列在所有排列中的字典序第几小

注意:
如果问比当前排列大/小 k的排列, 并不能先展开然后再放回.
!取模会丢失信息, 应该直接计算.
"""
# https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
# ! 取模的情况, len(s)很大


from typing import List

N = int(1e5) + 10
MOD = int(1e9 + 7)
fac = [1] * N
ifac = [1] * N
for i in range(1, N):
    fac[i] = fac[i - 1] * i % MOD
    ifac[i] = ifac[i - 1] * pow(i, MOD - 2, MOD) % MOD


def calRank(perm: List[int]) -> int:
    """求当前排列在所有排列中的字典序第几小(rank>=0)"""

    def add(i: int, val: int):
        while i <= n:
            bit[i] += val
            i += i & -i

    def preSum(i: int) -> int:
        res = 0
        while i > 0:
            res += bit[i]
            i &= i - 1
        return res

    n = len(perm)
    bit = [0] * (n + 1)
    for i in range(1, n + 1):
        add(i, 1)
    res = 0
    for i, v in enumerate(perm):
        res += preSum(v - 1) * fac[n - 1 - i] % MOD
        res %= MOD
        add(v, -1)
    return res  # 从0开始的排名


def calPerm(n: int, rank: int) -> List[int]:
    """求在1-n的所有排列中,字典序第几小(rank>=0)是谁"""
    fac = [1] * (n + 10)
    for i in range(1, n + 10):
        fac[i] = fac[i - 1] * i
    perm = [0] * n
    valid = [True] * (n + 1)
    for i in range(1, n + 1):
        order = rank // fac[n - i] + 1
        for j in range(1, n + 1):
            order -= valid[j]
            if order == 0:
                perm[i - 1] = j
                valid[j] = False
                break
        rank %= fac[n - i]
    return perm


if __name__ == "__main__":
    print(calRank([1, 2, 3, 4, 5, 6, 7, 8, 10, 9]))
    print(calPerm(10, 1))
