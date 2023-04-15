import sys

# !问 [0,2^D-1]内有多少个数与nums中的每个数按位与都不为0
# n<=20
# D<=60
# ai<=2^D
# !反向思考、按位考虑:与哪些数与运算会使得按位与为0,再用容斥原理计算
# 和n1+n2+..+nn -n1 and n2 -n1 and n3 -n1 and ... and nn-n1 +...
sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


n, d = map(int, input().split())
nums = list(map(int, input().split()))


def bit_count(n: int) -> int:
    return bin(n).count('1')


bad = 0
for state in range(1, (1 << n)):  # 枚举和哪些数与运算为0
    badDigit = 0
    for i in range(n):
        if (state >> i) & 1:
            badDigit |= nums[i]

    count = 2 ** (d - bit_count(badDigit))
    if bit_count(state) & 1:
        bad += count
    else:
        bad -= count

print(2 ** d - bad)

