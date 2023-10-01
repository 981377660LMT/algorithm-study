# 发糖果
# 每个人有一个讨厌的人pi(不是自己)
# 如果讨厌的人在他之前得到了糖果，那么他的生气值为ai
# 最小化生气值之和
# n<=2e5

# !https://atcoder.jp/contests/abc256/editorial/4135
# ! n个顶点n条边 Namori Graph(基环树森林)
# ! 内向基环树森林找环 每个环中必定有一个点作为起点 怒气值最小
# ! 只有在环里才能产生怒气 如果没环直接按照拓扑序发饼干 怒气值为0


import sys
from collections import defaultdict

from 基环树找到所有环 import cyclePartition


input = sys.stdin.readline
sys.setrecursionlimit(int(1e9))

n = int(input())
hates = [int(num) - 1 for num in input().split()]
scores = list(map(int, input().split()))

adjMap = defaultdict(set)
for i in range(n):
    adjMap[i].add(hates[i])

cycleGroup, *_ = cyclePartition(n, adjMap, directed=True)
res = 0
for group in cycleGroup:
    res += min(scores[i] for i in group)
print(res)
