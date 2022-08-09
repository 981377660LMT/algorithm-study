# n<=510
# len(words[i])<=8
# L<=100

# n个单词 每个单词都是01序列
# 从这些单词中选若干个按任意顺序连接成新单词，求长度为l的不同新单词的个数
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)

n, l = map(int, input().split())
words = []
for _ in range(n):
    words.append(input())


# TODO
