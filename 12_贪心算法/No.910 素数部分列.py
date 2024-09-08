# No.910 素数部分列
# https://yukicoder.me/problems/no/910
# 从一个由1，3，5，7，9组成的数字字符串中不断移除素数子序列，求可以操作的最大次数.
# 贪心.
# 3, 5, 7, 11, 19, 991
#
# !优先移除 3, 5, 7，然后移除 19, 991，最后移除 11

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    _ = int(input())
    s = input()

    nums = [int(v) for v in s]
    res = nums.count(3) + nums.count(5) + nums.count(7)  # 3, 5, 7
    c1 = 0
    c9 = 0
    for s in nums:
        if s == 1:
            c1 += 1
        elif s == 9:  # 19
            if c1:
                res += 1
                c1 -= 1
            else:
                c9 += 1
    while c9 >= 2 and c1:  # 991
        res += 1
        c9 -= 2
        c1 -= 1
    res += c1 // 2  # 11
    print(res)
