# Lyndon分解
# https://judge.yosupo.jp/problem/lyndon_factorization
# https://oi-wiki.org/string/lyndon/
# https://www.cnblogs.com/ptno/p/16418308.html

# !lyndon 串：一个字符串，如果他是他的最小后缀，那么他就是 lyndon 串。
# 还有一种定义是，在他的循环同构里他是字典序最小的那个。
# !lyndon 分解: 将一个字符串分解为若干个`字典序非严格递减`的 lyndon 串。

from typing import List


def lyndonFactorization(s: str) -> List[int]:
    """lyndon 分解."""
    n = len(s)
    res = [0]
    ptr = 0
    while ptr < n:
        i, j = ptr, ptr + 1
        while j < n and s[i] <= s[j]:
            if s[i] == s[j]:
                i += 1
                j += 1
            else:
                i = ptr
                j += 1
        len_ = j - i
        while ptr <= i:
            ptr += len_
            res.append(ptr)
    return res


if __name__ == "__main__":
    test = "babaabaab"
    splits = lyndonFactorization(test)
    for pre, cur in zip(splits, splits[1:]):
        print(test[pre:cur])

    s = str(input())
    for f in lyndonFactorization(s):
        print(f)
