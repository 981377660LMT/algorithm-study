# 给定一个长度为2*n的字符串s
# 判断是否存在一个长为n的字符串word
# !将word[::-1]插入到word中的某个位置(0<=pos<=n)后能够得到s
# 如果存在输出word和pos,否则输出-1
# n<=1e6

from typing import Tuple
from typing import Sequence


def findInsertPos(s: str) -> Tuple[Tuple[str, int], bool]:
    n = len(s) // 2
    ords = [ord(c) for c in s]
    MOD, BASE = 10**11 + 7, 1313131
    query1 = useStringHasher(ords, MOD, BASE)
    query2 = useStringHasher(ords[::-1], MOD, BASE)
    for pos in range(n + 1):  # 插入位置
        leftHash = query1(0, pos)
        rightHash = query1(pos + n, 2 * n)
        hash1 = (leftHash * pow(BASE, n - pos, MOD) + rightHash) % MOD  # 左右连接
        hash2 = query2(n - pos, 2 * n - pos)
        if hash1 == hash2:
            return (s[pos : pos + n][::-1], pos), True

    return ("", -1), False


def useStringHasher(ords: Sequence[int], mod=10**11 + 7, base=1313131):
    n = len(ords)
    prePow = [1] * (n + 1)
    preHash = [0] * (n + 1)
    for i in range(1, n + 1):
        prePow[i] = (prePow[i - 1] * base) % mod
        preHash[i] = (preHash[i - 1] * base + ords[i - 1]) % mod

    def sliceHash(left: int, right: int):
        """切片 `s[left:right]` 的哈希值"""
        if left >= right:
            return 0
        left += 1
        return (preHash[right] - preHash[left - 1] * prePow[right - left + 1]) % mod

    return sliceHash


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


if __name__ == "__main__":
    n = int(input())
    s = input()
    (word, insertPos), ok = findInsertPos(s)
    if not ok:
        print(-1)
        exit(0)
    print(word)
    print(insertPos)
