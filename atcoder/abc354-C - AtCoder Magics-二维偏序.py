# C - AtCoder Magics
# https://atcoder.jp/contests/abc354/tasks/abc354_c
# n个卡牌，有对应的强度ai和花费ci。对于两个卡牌i,j，如果ai>aj且ci<cj，则卡牌j会被丢弃。
# 不断进行如上操作，问最终的牌是哪些。

from typing import List, TypeVar

T = TypeVar("T")

INF = int(1e18)


def argSort(n: int, key) -> List[int]:
    """数组排序后的索引."""
    return sorted(range(n), key=key)


def rearrange(arr: List[T], order: List[int]) -> List[T]:
    """将数组按照order里的顺序重新排序."""
    res = arr[:]
    for i in range(len(order)):
        res[i] = arr[order[i]]
    return res


N = int(input())
A, C = [], []
for i in range(N):
    a, c = map(int, input().split())
    A.append(a)
    C.append(c)


pairs = [(A[i], C[i], i) for i in range(N)]
pairs.sort(key=lambda x: x[1])

res = []
preMax = -INF
for a, _, i in pairs:
    if a < preMax:
        continue
    preMax = a
    res.append(i)

res.sort()
print(len(res))
print(" ".join(str(x + 1) for x in res))
