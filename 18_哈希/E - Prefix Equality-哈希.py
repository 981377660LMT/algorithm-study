# 每个查询询问
# !nums1的前x项的集合是否与nums2的前y项的集合相等 (重复元素只算一次)
# n<=2e5
# numsi<=1e9
# x,y<=n


# 1.离线查询(排序)+双指针 莫队的思想
# 固定x之后 就可以尺取寻找y的边界

# !2. 异或哈希/异或前缀和 用随机数产生哈希值 用异或来计算区间所含集合的哈希值
# Zobrist Hash
# !xorではなく和を使うと個数に対応したハッシュが作れる


import sys
import os
from typing import List, Tuple
from collections import defaultdict
from itertools import accumulate
from operator import xor
from random import randint

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def genHash(nums1: List[int], nums2: List[int]) -> Tuple[List[int], List[int]]:
    """随机数+异或来生成前缀的集合的哈希值"""
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


def main() -> None:
    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))
    q = int(input())
    Q = []
    for _ in range(q):
        x, y = map(int, input().split())
        Q.append((x, y))
    res = [0] * q

    # !检验两次 防止哈希冲突
    for _ in range(2):
        pre1, pre2 = genHash(nums1, nums2)
        for i, (x, y) in enumerate(Q):
            if pre1[x] == pre2[y]:
                res[i] += 1

    for v in res:
        if v == 2:
            print("Yes")
        else:
            print("No")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()


# 同样的应用:
# ! q次查询子串s[l:r] 能否经过重排形成回文串
# 至多有一种字符出现奇数次
# 给每种字符赋互异的随机权值，然后对字符串转化成的权值串进行异或前缀和
# 只需要判断区间异或和是否为 0 或者某个字符的随机权值即可。
