# 1<=L,R<=1e18
# 每个数num在 黑板上写num次 求最终的长度
# !前缀和相减 按照位数分类计算

import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


L, R = map(int, input().split())


def cal(upper: int) -> int:
    """[1, upper]内的答案"""
    res = 0
    for i in range(20):
        left, right = 10 ** i, 10 ** (i + 1) - 1
        wordLen = i + 1
        if right >= upper:
            count = upper - left + 1
            res += wordLen * (left + upper) * count // 2
            break
        count = right - left + 1
        res += wordLen * (left + right) * count // 2
    return res


print((cal(R) - cal(L - 1)) % MOD)
