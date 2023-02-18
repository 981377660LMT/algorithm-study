# multiset哈希
# Zobrist HashingでXOR使うところを、
# !和を使うことで個数にも対応するテクニック、
# 普通に超便利なので将棋AI以外にも浸透したほうが良いと思ってるんだけど、
# 普通に浸透してるのかな。
# https://blog.hamayanhamayan.com/entry/2017/05/24/154618

from collections import defaultdict
from random import randint
from typing import List


pool = defaultdict(lambda: randint(1, (1 << 63) - 1))


def mutisetHash(mutiset: List[int]) -> int:
    """多重集合哈希值"""
    return sum(pool[x] for x in mutiset)


nums1, nums2 = [1, 2, 3, 4, 5], [1, 2, 3, 4, 5]
print(mutisetHash(nums1) == mutisetHash(nums2))  # True
