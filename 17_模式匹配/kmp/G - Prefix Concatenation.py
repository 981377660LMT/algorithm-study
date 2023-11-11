"""t至少需要用多少个s的前缀连接成

|s,t|<=5e5
两段前缀 变为 一段前缀+一段后缀
"""

import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def getNext(needle: str) -> List[int]:
    """kmp O(n)求 `needle`串的 `next`数组

    `next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度
    """
    next = [0] * len(needle)
    j = 0

    for i in range(1, len(needle)):
        while j and needle[i] != needle[j]:  # 1. fallback后前进：匹配不成功j往右走
            j = next[j - 1]

        if needle[i] == needle[j]:  # 2. 匹配：匹配成功j往右走一步
            j += 1

        next[i] = j

    return next


def main() -> None:
    s = input()
    t = input()
    n = len(t)
    next_ = getNext(s + "#" + t)[-n:]  # 两段前缀 变为 一段前缀+一段后缀

    res, pos = 0, n - 1
    while pos > 0 and next_[pos - 1] > 0:  # lcp长度>1
        pos -= next_[pos - 1]
        res += 1
    if pos == 0:
        print(res)
    else:
        print(-1)


if sys.argv[0].startswith(r"e:\test"):
    while True:
        try:
            main()
        except (EOFError, ValueError):
            break
else:
    main()
