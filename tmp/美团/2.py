# TODO tcp udp http 端口

# 模拟+枚举候选人

from collections import Counter
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(1e18)

n = int(input())
types1 = list(map(int, input().split()))
types2 = list(map(int, input().split()))

# 注意到数字必须要大于一半 因此候选人数量很少
# 枚举这些可能合法的候选人
counter1 = Counter(types1)
counter2 = Counter(types2)
counter3 = counter1 + counter2
half = (n + 1) // 2
cands = [k for k, v in counter3.items() if v >= half]
if not cands:
    print(-1)
    exit(0)


def check(cand: int) -> int:
    """正面朝上为cand时需要反转的魔法数量"""
    remain = half
    visited = [False] * n

    for i, num in enumerate(types1):
        if num == cand:
            visited[i] = True
            remain -= 1
            if remain == 0:
                return 0

    res = 0
    for i, num in enumerate(types2):
        if visited[i]:
            continue
        if num == cand:
            res += 1
            remain -= 1
            if remain == 0:
                return res

    return INF


res = [check(num) for num in cands]
min_ = min(res)
print(min_ if min_ <= n else -1)
