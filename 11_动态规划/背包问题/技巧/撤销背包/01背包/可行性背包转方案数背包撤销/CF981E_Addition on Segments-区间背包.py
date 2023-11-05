# https://www.luogu.com.cn/problem/CF981E

# 给一个长度为n的序列(初始全为0)和q条操作(以(left,right,x)表示将[left,right]中的每个数都加上x(1<=x<=n)
# 对于1≤k≤n,求哪些k满足:选出若干条操作后序列最大值为k.(对于一个k,每条操作至多用一次)
# n,q<=1e4，区间背包
#
# !扫描线+撤销背包
# !时间复杂度O(nk)


from typing import List, Tuple
from Knapsack01Removable import Knapsack01Removable

MOD = int(1e9 + 7)  # 大素数


def additionOnSegments(n: int, operations: List[Tuple[int, int, int]]) -> List[int]:
    ins = [[] for _ in range(n + 1)]
    outs = [[] for _ in range(n + 1)]
    for start, end, x in operations:
        ins[start].append(x)
        outs[end].append(x)

    ok = [False] * (n + 1)
    dp = Knapsack01Removable(n, MOD)

    for i in range(1, n + 1):
        for w in ins[i]:
            dp.add(w)
        for j in range(1, n + 1):
            ok[j] |= dp.query(j) > 0
        for w in outs[i]:
            dp.remove(w)
    return [i for i in range(1, n + 1) if ok[i]]


if __name__ == "__main__":
    import sys

    input = sys.stdin.readline

    n, q = map(int, input().split())
    op = []
    for _ in range(q):
        left, right, x = map(int, input().split())
        op.append((left, right, x))
    res = additionOnSegments(n, op)
    print(len(res))
    print(*res)
