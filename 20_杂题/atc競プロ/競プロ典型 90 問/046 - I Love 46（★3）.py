# 三个数组 问能有多少个三元组 和为46的倍数
# 同じ意味のものをまとめて考える

from collections import Counter
import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n = int(input())
counter1 = Counter(num % 46 for num in map(int, input().split()))
counter2 = Counter(num % 46 for num in map(int, input().split()))
counter3 = Counter(num % 46 for num in map(int, input().split()))

res = 0
for k1, v1 in counter1.items():
    for k2, v2 in counter2.items():
        for k3, v3 in counter3.items():
            if (k1 + k2 + k3) % 46 == 0:
                res += v1 * v2 * v3
print(res)
