# n<=12
# Q<=50
# wi<2^60
# 有一些限制[i1,i2,i3,expected]
# !求数组nums的个数 使得nums[i1] or nums[i2] or nums[i3] 为 expected
# !按位考虑 二进制枚举 计算每位的贡献
# !选取数字相当于 对每一位讨论 相乘

# !时间复杂度 2^n*logwi*(n+q)
import sys


sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n, q = map(int, input().split())
queries = []
for _ in range(q):
    i1, i2, i3, expected = map(int, input().split())
    i1, i2, i3 = i1 - 1, i2 - 1, i3 - 1
    queries.append((i1, i2, i3, expected))


def cal(bit: int) -> int:
    res = 1 << n  # 每一个数在这个位上取0/1
    for state in range(1 << n):
        for i1, i2, i3, expected in queries:
            if (state >> i1 & 1) | (state >> i2 & 1) | (state >> i3 & 1) != (expected >> bit & 1):
                res -= 1  # 不合法的state
                break
    return res


res = 1
for bit in range(60):
    res *= cal(bit)
    res %= MOD
print(res)
