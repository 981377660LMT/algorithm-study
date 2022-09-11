# 有一把魔法权杖，权杖上有 n 颗并排的法术石(编号为 1 到 n)。
# 每颗法术石具有一个能量值，权杖的法术强度等同于法术石的最小能量值。
# 权杖可以强化，一次强化可以将`两颗相邻`的法术石融合为一颗，
# 融合后的能量值为这两颗法术石能量值之和。现在有 m 次强化的机会，
# 请问权杖能 强化到的最大法术强度是多少？
# !n,m<=1e5

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


n, m = map(int, input().split())
nums = list(map(int, input().split()))


def check(mid: int) -> bool:
    """强化到的最大法术强度是否能达到 mid"""

    # 相邻元素合并(消除)用栈
    ...


left, right = min(nums), sum(nums)
while left <= right:
    mid = (left + right) // 2
    if check(mid):
        left = mid + 1
    else:
        right = mid - 1

print(right)
