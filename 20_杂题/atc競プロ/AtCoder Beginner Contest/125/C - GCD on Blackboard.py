# 将一个数变为1到1e9的任一个数
# !求变化后数组gcd最大值
# n<=1e5 1<=ai<=1e9

# !枚举变化哪个位置的数+前后缀分解


from itertools import accumulate
from math import gcd
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    preGcd = [0] + list(accumulate(nums, gcd))
    sufGcd = ([0] + list(accumulate(nums[::-1], gcd)))[::-1]
    print(max(gcd(preGcd[i], sufGcd[i + 1]) for i in range(n)))
