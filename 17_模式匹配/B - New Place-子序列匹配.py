"""
给定两个长度为n的字符串s和t
每次可以将s的第一个字符删除,然后插入到s的任意位置
问最少需要多少次操作可以使得s和t相等
如果无法使得s和t相等,输出-1

子序列匹配
"""
from collections import Counter


def newPlace(s: str, t: str) -> int:
    counterS, counterT = Counter(s), Counter(t)
    if counterS != counterT:
        return -1

    s = s[::-1]
    t = t[::-1]
    hit = 0
    for i in range(n):
        if s[hit] == t[i]:
            hit += 1
    return n - hit


if __name__ == "__main__":
    n = int(input())
    s = input()
    t = input()
    print(newPlace(s, t))
