"""
给定n个字符串,对于每个字符串 si,问 maxLCP(si,sj)i≠j
其中LCP是最长公共前缀。

最大的最长公共前缀一定在字典序上前后两个的字符串之间
因此将这n个字符串按字典序排序,求每个字符串与其相邻(也可以计算与后面那个)的字符串的LCP,取最大值即可
"""


from typing import List


def karuta(words: List[str]) -> List[int]:
    n = len(words)
    wordsWithIndex = [(i, word) for i, word in enumerate(words)]
    wordsWithIndex.sort(key=lambda x: x[1])
    res = [0] * n
    for (i1, s1), (i2, s2) in zip(wordsWithIndex, wordsWithIndex[1:]):
        lcp = calLCP(s1, s2)
        res[i1] = max(res[i1], lcp)
        res[i2] = max(res[i2], lcp)
    return res


def calLCP(s1: str, s2: str) -> int:
    count = 0
    for a, b in zip(s1, s2):
        if a == b:
            count += 1
        else:
            break
    return count


if __name__ == "__main__":
    n = int(input())
    words = [input() for _ in range(n)]
    print(*karuta(words), sep="\n")
