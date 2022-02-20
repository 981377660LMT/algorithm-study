# 沿着长 N 米、宽 1 米的走廊，连续铺有 N 块边长为 1 米的正方形地砖。
# 假设所有颜色分为 26 种，用小写字母 a 到 z 表示，
# 给定 N 块地砖的初始颜色， 每次可以选择一种颜色，
# 然后将最多M块连续的地砖染成该颜色。那么至少要进行多少次染色，才能将所有地砖染成同一颜色?

# 枚举+贪心
# 1 ≤ M ≤ N ≤ 104

# 总结:要多用python的next函数，可以非常方便地处理边界

from collections import defaultdict
from math import ceil


N, M = map(int, input().split())
string = input()
indexes = defaultdict(set)
for i, char in enumerate(string):
    indexes[char].add(i)


res = ceil(N / M)
for color in indexes.keys():
    curOpt = 0
    need = set(range(N)) - indexes[color]
    start = next((i for i in range(N) if i in need), N)
    while start < N:
        curOpt += 1
        start += M
        start = next((i for i in range(start, N) if i in need), N)

    res = min(res, curOpt)

print(res)
