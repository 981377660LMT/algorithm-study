# 我们希望知道一个 N 位字符串 S 的前缀是否具有循环节。
# 换言之，对于每一个从头开始的长度为 i（i>1）的前缀，是否由重复出现的子串 A 组成，即 AAA…A （A 重复出现 K 次,K>1）。
# 如果存在，请找出最短的循环节对应的 K 值（也就是这个前缀串的所有可能重复节中，最大的 K 值）。

import sys
from typing import List

# input = lambda: sys.stdin.readline().strip()


def getNext(pattren: str) -> List[int]:
    next = [0] * len(pattren)
    j = 0

    for i in range(1, len(pattren)):
        while j and pattren[i] != pattren[j]:
            j = next[j - 1]

        if pattren[i] == pattren[j]:
            j += 1

        next[i] = j

    return next


count = 1
while True:
    n = int(input())
    if n == 0:
        break

    print(f"Test case #{count}")
    count += 1

    # next[i]表示[:i+1]这一段字符串中最长公共前后缀(不是原串)的长度
    next = getNext(input())

    for i in range(1, n):
        t = (i + 1) - next[i]
        if t and (i + 1) > t and (i + 1) % t == 0:
            print((i + 1), (i + 1) // t)
    print()

