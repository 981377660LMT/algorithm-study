# HandStand (倒立)

# !k次翻转区间flip后连续1的最大长度
# 01串区间翻转k次


# Solution
# 1.注意翻转只考虑翻转一段连续0区间
# !2.groupby后转化为fix模型

from itertools import groupby
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, k = map(int, input().split())
    s = input()  # !01串

    groups = [(char, len(list(group))) for char, group in groupby(s)]
    res, left, m = 0, 0, len(groups)
    count, remain = 0, k
    for right in range(m):
        char, size = groups[right]
        count += size
        if char == "0":
            remain -= 1
        while left <= right and remain < 0:
            count -= groups[left][1]
            if groups[left][0] == "0":
                remain += 1
            left += 1
        res = max(res, count)

    print(res)
