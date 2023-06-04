# Prefix-Substring LCS
# s的前缀与t的子串的最长公共子序列(LCS)
# https://judge.yosupo.jp/problem/prefix_substring_lcs


from typing import Any, Sequence


def max(a: int, b: int) -> int:
    return a if a > b else b


def min(a: int, b: int) -> int:
    return a if a < b else b


class PrefixSubstringLCS:
    __slots__ = "_dp1"

    def __init__(self, seq1: Sequence[Any], seq2: Sequence[Any]) -> None:
        n1, n2 = len(seq1), len(seq2)
        dp1, dp2 = [[0] * (n2 + 1) for _ in range(n1 + 1)], [[0] * (n2 + 1) for _ in range(n1 + 1)]
        dp1[0] = list(range(n2 + 1))
        for i in range(1, n1 + 1):
            cache1, cache2, cache3 = dp1[i], dp2[i], dp1[i - 1]
            for j in range(1, n2 + 1):
                if seq1[i - 1] == seq2[j - 1]:
                    cache1[j] = cache2[j - 1]
                    cache2[j] = cache3[j]
                else:
                    max_, min_ = cache3[j], cache2[j - 1]
                    if max_ < min_:
                        max_, min_ = min_, max_
                    cache1[j] = max_
                    cache2[j] = min_
        self._dp1 = dp1

    def query(self, a: int, b: int, c: int) -> int:
        """查询 seq1[:a) 和 seq2[b:c) 的最长公共子序列长度."""
        cache = self._dp1[a]
        return sum(cache[i] <= b for i in range(b + 1, c + 1))


if __name__ == "__main__":
    q = int(input())
    s = input()
    t = input()
    ords1, ords2 = [ord(ch) for ch in s], [ord(ch) for ch in t]
    lcs = PrefixSubstringLCS(ords1, ords2)
    for _ in range(q):
        a, b, c = map(int, input().split())
        print(lcs.query(a, b, c))
