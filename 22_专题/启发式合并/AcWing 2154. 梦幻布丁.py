'''
思路
用集合来维护每种颜色的位置信息，变颜色相当于把两个集合合并成一个集合，用启发式合并将
X集合合并到Y集合前，遍历X中所有位置，如果某位置pos前一个位置或者后一个位置的颜色是Y，
那都会导致总颜色分段计数减去1，合并集合同时维护住总的分段计数即可
https://www.acwing.com/activity/content/code/content/1100516/
'''
# n,m<=1e5
# colors数<1e6

from collections import defaultdict
from itertools import groupby


n, m = map(int, input().split())
nums = list(map(int, input().split()))

colorMap = defaultdict(set)
groupCount, index = 0, 0
for color, group in groupby(nums):
    groupCount += 1
    for _ in group:
        colorMap[color].add(index)
        index += 1


for _ in range(m):
    opt, *rest = list(map(int, input().split()))
    if opt == 1:  # 若第一个数是 1，则表示好友要将某种颜色🍮变成另一种颜色
        x, y = rest
        if x == y:
            continue

        set1, set2 = colorMap[x], colorMap[y]
        # 当前位置前后只要有一样的颜色，颜色的分段就会减少
        for pos in set1:
            if pos - 1 in set2:
                groupCount -= 1
            if pos + 1 in set2:
                groupCount -= 1

        # 启发式合并
        set1, set2 = sorted([set1, set2], key=len)
        set2 |= set1
        colorMap[y] = set2
        colorMap[x] = set()

    elif opt == 2:  # 若第一个数是 2，则表示好友要询问目前有多少段颜色(group)
        print(groupCount)
