# 选卡片 ban里的卡片不能同时选
# 两个人最后选的卡片分数和相等 但是卡片的集合不同
# 输出两人的卡片集合(一种即可)
# n<=88
# 0<=sum(cards)<=8888

# !鸽巢原理:因为和只有8889种
# !因此枚举出8890种不同的卡片集合 必定有一对和相等的卡片集合
# !在此基础上回溯搜索 8890*88 次(只要搜8890次就会找到两个一样的 每次搜索到尽头为n)

from collections import defaultdict
import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


n, q = map(int, input().split())
cards = list(map(int, input().split()))
ban = defaultdict(set)
for _ in range(q):
    a, b = map(int, input().split())
    a, b = a - 1, b - 1
    ban[a].add(b)
    ban[b].add(a)


def logFatal(arr1: List[bool], arr2: List[bool]) -> None:
    res1 = [i + 1 for i, v in enumerate(arr1) if v]
    res2 = [i + 1 for i, v in enumerate(arr2) if v]
    print(len(res1))
    print(*res1)
    print(len(res2))
    print(*res2)
    exit(0)


def bt(index: int, curSum: int, selected: List[bool]) -> None:
    if index == n:
        if visited[curSum]:
            logFatal(visited[curSum], selected)
        else:
            visited[curSum] = selected[:]
        return

    # 不选
    bt(index + 1, curSum, selected)

    # 选(不冲突的情况下)
    if all(not selected[i] for i in ban[index]):
        selected[index] = True
        bt(index + 1, curSum + cards[index], selected)
        selected[index] = False


visited = defaultdict(list)
bt(0, 0, [False] * n)

