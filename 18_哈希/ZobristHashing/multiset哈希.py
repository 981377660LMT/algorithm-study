# multiset哈希
# Zobrist HashingでXOR使うところを、
# !和を使うことで個数にも対応するテクニック、
# 普通に超便利なので将棋AI以外にも浸透したほうが良いと思ってるんだけど、
# 普通に浸透してるのかな。
# https://blog.hamayanhamayan.com/entry/2017/05/24/154618

from collections import Counter, defaultdict
from random import randint
from typing import List, Set


pool = defaultdict(lambda: randint(1, (1 << 61) - 1))


def mutisetHash(mutiset: List[int]) -> int:
    """多重集合哈希值"""
    return sum(pool[x] for x in mutiset)


def setHash(s: Set[int]) -> int:
    """集合哈希值"""
    hash_ = 0
    for x in s:
        hash_ ^= pool[x]
    return hash_


nums1, nums2 = [1, 2, 3, 4, 5], [1, 2, 3, 4, 5]
print(mutisetHash(nums1) == mutisetHash(nums2))  # True
hash1, hash2 = setHash(set(nums1)), setHash(set(nums2))
print(hash1 == hash2)  # True
# check with bruteforce
for _ in range(10000):
    s1, s2 = set(), set()
    for _ in range(3):
        s1.add(randint(1, 10))
        s2.add(randint(1, 10))
    assert (s1 == s2) == (setHash(s1) == setHash(s2))

    counter1, counter2 = Counter(), Counter()
    for _ in range(3):
        counter1[randint(1, 10)] += 1
        counter2[randint(1, 10)] += 1

    assert (counter1 == counter2) == (
        mutisetHash(list(counter1.elements())) == mutisetHash(list(counter2.elements()))
    )
