# 1. 15Puzzle Game(华容道???) 是否存在解
# https://cp-algorithms.com/others/15-puzzle.html


from typing import List


def isExisted(perm: List[int]) -> bool:
    inv = 0
    for i in range(16):
        if perm[i]:
            for j in range(i):
                if perm[j] > perm[i]:
                    inv += 1
    for i in range(16):
        if perm[i] == 0:
            inv += 1 + i // 4
    return inv & 1 == 0


print(isExisted([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0]))
