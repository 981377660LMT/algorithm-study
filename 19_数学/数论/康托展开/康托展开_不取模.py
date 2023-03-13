"""
康托展开 - 有重复元素
求字典序第k小的排列/当前排列在所有排列中的字典序第几小
"""

# ! 不取模的情况, len(s)很小

from collections import Counter
from functools import lru_cache
from typing import List, Sequence, TypeVar


T = TypeVar("T", str, int)


@lru_cache(None)
def fac(n: int) -> int:
    if n <= 1:
        return 1
    return n * fac(n - 1)


def calRank(s: Sequence[T]) -> int:
    """求当前排列在所有排列中的字典序第几小(rank>=0)"""
    n = len(s)
    counter = Counter(s)
    keys = sorted(counter)
    res = 0
    for i, char in enumerate(s):
        suf = fac(n - i - 1)
        for count in counter.values():
            suf //= fac(count)  # !后面位置的组合数
        for smaller in keys:
            if smaller >= char:
                break
            res += counter[smaller] * suf
        counter[char] -= 1
    return res


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
            suf = fac(n - i - 1)
            for count in counter.values():
                suf //= fac(count)  # !后面位置的组合数
            if suf > rank:
                res.append(char)
                break
            else:
                rank -= suf
                counter[char] += 1
    return res


if __name__ == "__main__":

    # https://yukicoder.me/problems/no/1311
    # !求出1-n的排列中,字典序第k小的排列的`逆置换`的名次
    k, n = map(int, input().split())
    s = calPerm(list(range(1, n + 1)), k)
    rs = [s.index(i) + 1 for i in range(1, n + 1)]
    print(calRank(rs))
