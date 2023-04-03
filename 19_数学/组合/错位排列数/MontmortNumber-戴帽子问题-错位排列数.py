# https://judge.yosupo.jp/problem/montmort_number_mod
# Montmort配对(MontmortNumber)
# 戴帽子问题
# n<=1e6
# 参加集会的n个人将他们的帽子放在一起，
# 会后每人任取一顶帽子戴上。求恰好有k个人戴对自己的帽子的方案数
# !扰乱排列(错位排列)：使得元素中没有一个出现在其原始位置


from typing import List


def derangement(n: int, mod=None) -> List[int]:
    """对k=0,1,...,N,求错位排列数."""
    if n < 0:
        return []
    elif n == 0:
        return [0]
    elif n == 1:
        return [0, 0]
    elif mod == 1:
        return [0] * (n + 1)

    res = [0] * (n + 1)
    res[2] = 1
    b, c = 0, 1

    for k in range(3, n + 1):
        b, c = c, (k - 1) * (b + c)
        if mod is not None:
            c %= mod
        res[k] = c

    return res


n, MOD = map(int, input().split())
res = derangement(n, MOD)
print(*res[1:])
