# 值班表
# 给定高桥的n天值班情况s。
# n<=2e5
# 问满足下述条件的青木的n天值班情况t数量
# 满足每天他俩至少有一人值班，且青木的值班情况是关于m的循环，其中 m|n,m<n。

# F - Shift Table-约数系容斥原理、约数系包除原理

from typing import List

MOD = 998244353


def shiftTable(work: List[bool]) -> int:
    n = len(work)
    res = [0] * n
    for fac in range(1, n):
        if not n % fac:
            continue
        cur = 1
        for i in range(n):
            if all(work[k] for k in range(i, n, fac)):
                cur *= 2  # 这一天可以休息
                cur %= MOD
        res[fac] += cur

        # !减去重复计算的
        for j in range(2 * fac, n, fac):
            if not n % j:
                continue
            res[fac] -= res[j]
            res[fac] %= MOD

    return sum(res) % MOD


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    s = input()
    work = [True if v == "#" else False for v in s]
    print(shiftTable(work))
