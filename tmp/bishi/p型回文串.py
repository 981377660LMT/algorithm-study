# 小明和小红都很喜欢收集字符串。小明喜欢收集周期为P的字符串，
# 小红喜欢收集回文串。当一个回文串的周期为P时，我们称它为“P型回文串”。
# 现在给你一个长度为P的整数倍的字符串，请问你最少改变几个字符能将它变成一个P型回文串?
# 一个字符串的周期为Р，当且仅当字符串中任意两个距离之差为P的位置上的字符均相等。
# 一个字符串是回文串，当且仅当字符串前后翻转后与原串相等。
# p型回文串

from collections import defaultdict


def solve(n: int, p: int, s: str) -> int:
    res = 0
    group = [[] for _ in range(p)]
    for i in range(n):
        group[i % p].append(s[i])

    for i in range(p // 2):
        count = 0
        counter = defaultdict(int)
        for char in group[i]:
            counter[char] += 1
            count += 1
        if p - i - 1 > i:
            for char in group[p - i - 1]:
                counter[char] += 1
                count += 1
        max_ = max(counter.values())
        res += count - max_
    return res


t = int(input())
for _ in range(t):

    n, p = map(int, input().split())
    s = input()
    print(solve(n, p, s))


# cdccdc  2
# zpzzpz 0
