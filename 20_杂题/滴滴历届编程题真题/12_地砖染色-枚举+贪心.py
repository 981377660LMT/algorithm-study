# 地砖染色
# 沿着长 N 米、宽 1 米的走廊，连续铺有 N 块边长为 1 米的正方形地砖。
# 假设所有颜色分为 26 种，用小写字母 a 到 z 表示，
# 给定 N 块地砖的初始颜色， 每次可以选择一种颜色，
# !然后将最多 M 块连续的地砖染成该颜色。
# 那么至少要进行多少次染色，才能将所有地砖染成同一颜色?

# 枚举最后染成的颜色
# 1 ≤ M ≤ N ≤ 1e4


from math import ceil
import string


n, m = map(int, input().split())
s = input()


res = ceil(n / m)
for char in string.ascii_lowercase:
    cand = 0
    index = 0
    while index < n:
        if s[index] != char:
            cand += 1
            index += m
        else:
            index += 1
    res = min(res, cand)

print(res)
