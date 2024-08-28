from typing import List
from kmp import getNext


INF = int(1e18)


def compressString(pre: str, post: str) -> int:
    """pre的后缀和post的前缀的最大公共长度"""
    cat = post + "#" + pre
    next_ = getNext(cat)
    return next_[-1]


def compressNums(pre: List[int], post: List[int]) -> int:
    """pre的后缀和post的前缀的最大公共长度"""
    cat = post + [INF] + pre
    next_ = getNext(cat)
    return next_[-1]


def compressStringNaive(pre: str, post: str) -> int:
    """pre的后缀和post的前缀的最大公共长度"""
    n = len(pre)
    m = len(post)
    for res in range(min(n, m), -1, -1):
        if pre[n - res :] == post[:res]:
            return res
    return 0


def compressNumsNaive(pre: List[int], post: List[int]) -> int:
    """pre的后缀和post的前缀的最大公共长度"""
    n = len(pre)
    m = len(post)
    for res in range(min(n, m), -1, -1):
        if pre[n - res :] == post[:res]:
            return res
    return 0


if __name__ == "__main__":
    print(compressString("abab", "ababab"))
    print(compressNums([1, 2, 3], [2, 3, 4]))
    print(compressStringNaive("abab", "ababab"))
    print(compressNumsNaive([1, 2, 3], [2, 3, 4]))
