# 翻转前半与后半 => 懒更新

import sys
from typing import List, Tuple

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# T=1 のとき : S の A 文字目と B 文字目を入れ替える
# T=2 のとき : S の前半 N 文字と後半 N 文字を入れ替える


def ipfl(s: str, operations: List[Tuple[int, int, int]]) -> str:
    n = len(s) // 2
    sb = list(s)
    isRev = False
    for kind, i, j in operations:
        if kind == 1:  # !交换i,j
            if isRev:
                i = (i + n) if i < n else (i - n)
                j = (j + n) if j < n else (j - n)
            sb[i], sb[j] = sb[j], sb[i]
        else:  # !前半段与后半段交换
            isRev = not isRev

    if isRev:
        sb = sb[n:] + sb[:n]
    return "".join(sb)


if __name__ == "__main__":
    n = int(input())
    s = input()
    q = int(input())
    operations = []
    for _ in range(q):
        kind, i, j = map(int, input().split())
        i, j = i - 1, j - 1
        operations.append((kind, i, j))
    res = ipfl(s, operations)
    print(res)
