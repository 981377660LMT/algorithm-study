"""
康托展开 - 有重复元素
求字典序第k小的排列/当前排列在所有排列中的字典序第几小
"""

# ! 注意不取模的情况

from collections import Counter
from typing import List, Sequence, TypeVar


MOD = int(1e9 + 7)
fac = [1]
ifac = [1]
for i in range(1, int(2e5) + 10):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


T = TypeVar("T", str, int)


def calRank(s: Sequence[T]) -> int:
    """求当前排列在所有排列中的字典序第几小(rank>=0)"""
    n = len(s)
    counter = Counter(s)
    keys = sorted(counter)

    res = 0
    for i, char in enumerate(s):
        suf = fac[n - i - 1]

        for count in counter.values():
            suf *= ifac[count]  # !后面位置的组合数
            suf %= MOD

        for smaller in keys:
            if smaller >= char:
                break
            res += counter[smaller] * suf
            res %= MOD

        counter[char] -= 1

    return res % MOD


def calPerm(s: Sequence[T], rank: int) -> List[T]:
    """求在所有排列中,字典序第几小(rank>=0)是谁"""
    n = len(s)
    counter = Counter(s)
    keys = sorted(counter)

    res = []
    for i in range(n):
        for char in keys:
            if counter[char] == 0:
                continue

            counter[char] -= 1
            suf = fac[n - i - 1]  # factorial(n-i-1)

            for count in counter.values():
                suf *= ifac[count]  # ! factorial(count)
                suf %= MOD

            if suf > rank:
                res.append(char)
                break
            else:
                rank -= suf
                counter[char] += 1

    return res


if __name__ == "__main__":
    assert calRank("cba") == 5
    assert calRank("abc") == 0
    assert calRank([3, 4, 1, 5, 2]) == 61

    assert calPerm("cba", 1) == ["a", "c", "b"]
    assert calPerm("aab", 1) == ["a", "b", "a"]
