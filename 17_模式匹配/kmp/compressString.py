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


if __name__ == "__main__":
    print(compressString("abab", "ababab"))
    print(compressNums([1, 2, 3], [2, 3, 4]))
