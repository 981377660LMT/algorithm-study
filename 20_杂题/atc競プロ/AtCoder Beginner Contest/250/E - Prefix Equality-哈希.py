# 每个查询询问
# !nums1的前x项的集合是否与nums2的前y项的集合相等 (重复元素只算一次)
# n<=2e5
# numsi<=1e9
# x,y<=n


# !1. 哈希 用随机数产生哈希值 用异或来计算区间所含集合的哈希值


from collections import defaultdict
from itertools import accumulate
from operator import xor
from random import randint
import sys
from typing import List, Tuple


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def genHash(nums1: List[int], nums2: List[int]) -> Tuple[List[int], List[int]]:
    """随机数+异或来产生区间前缀集合的哈希值"""
    pool = defaultdict(lambda: randint(1, (1 << 63) - 1))
    res1, visited1 = [0] * len(nums1), set()
    res2, visited2 = [0] * len(nums2), set()
    for i, (a, b) in enumerate(zip(nums1, nums2)):
        if a not in visited1:  # !はじめて出るときだけxorをとる
            visited1.add(a)
            res1[i] = pool[a]
        if b not in visited2:
            visited2.add(b)
            res2[i] = pool[b]
    return ([0] + list(accumulate(res1, xor)), [0] + list(accumulate(res2, xor)))


n = int(input())
nums1 = list(map(int, input().split()))
nums2 = list(map(int, input().split()))
q = int(input())
Q = []
for _ in range(q):
    x, y = map(int, input().split())
    Q.append((x, y))

res = [0] * q
# !防止哈希冲突
for _ in range(2):
    preHash1, preHash2 = genHash(nums1, nums2)
    for i, (x, y) in enumerate(Q):
        if preHash1[x] == preHash2[y]:
            res[i] += 1

for v in res:
    if v == 2:
        print("Yes")
    else:
        print("No")
