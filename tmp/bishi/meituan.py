################################################
# !把第i个元素移到最左边
# 可以链表/set
# from sortedcontainers import SortedList

# n, m = map(int, input().split())
# indexes = [int(num) - 1 for num in input().split()]

# sl = SortedList()
# for i in range(1, n + 1):
#     sl.add((0, i))

# id = -1
# for i in indexes:
#     cur = sl.pop(i)
#     sl.add((id, cur[1]))
#     id -= 1  # 插入到最左边(指定插入位置)

# for i in range(n):
#     print(sl[i][1], end=" ")

# !把编号为i的元素移到最左边
# 可以set
# n, m = map(int, input().split())
# nums = list(map(int, input().split()))
# mp = {i: (0, i) for i in range(1, n + 1)}
# sl = SortedList()
# for i in range(1, n + 1):
#     sl.add((0, i))

# id = -1
# for num in nums:
#     pre = mp[num]
#     index = sl.index(pre)
#     sl.pop(index)
#     sl.add((id, pre[1]))
#     mp[num] = (id, pre[1])
#     id -= 1
# for i in range(n):
#     print(sl[i][1], end=" ")
################################################
# 一个串能否划分成目标的m个子串
# !注意有相同元素的场合
# from collections import Counter
# from itertools import combinations


# n, m = map(int, input().split())
# big = input()  # 原串
# lens = list(map(int, input().split()))
# small = []  # 各个子串
# for _ in range(m):
#     small.append(input())
# target = Counter(small)

# if m == 1:
#     print(1 if big == small[0] else 0)
# else:
#     # 暴力选出m-1个分割点的位置
#     res = 0
#     for splits in combinations(range(1, n), m - 1):
#         parts = []
#         splits = (0,) + splits + (n,)
#         for pre, cur in zip(splits, splits[1:]):
#             parts.append(big[pre:cur])
#         if Counter(parts) == target:
#             res += 1
#     print(res)
################################################
# 堆+反悔
################################################
# 美团笔试t4
# 机器人收衣服
# 小团正忙看用机器人收衣服!因为快要下雨了，小团找来了不少机器人帮忙收衣服。
# 他有n件衣版从左到右成一行排列，所在位置分别为1~n,
# 在每个位置上已经有一个就绪的机器人可以帮忙收衣服，但第i个位置上的机器人需要pi的电量来启动。
# 然后这个机器人会用ti的时间收衣服，当它收完当前衣服后，
# 会尝试去收紧邻的右边的一件衣服(如果存在的话)，即i+1处的衣服，
# 如果i+1处的衣服已经被其他机器人收了或者其他机器人正在收，这个机器人就会进入休眠状态，不再收衣服。
# 不过如果机器人没有休眠，它会同样以ti时间来收这件1+1处的衣服(注意，不是t+1的时间，收衣服的时间为每个机器人固有属性)，
# 然后它会做同样的检测来看能否继续收i+2处的衣服，一直直到它进入休眠状态或者右边没有衣服可以收了。
# 形象地来说，机器人会一直尝试往右边收衣服，收k件的话就耗费k*ti;的时间,
# 但是当它遇见其他机器人工作的痕迹，就会认为后面的事情它不用管了,开始摸鱼，进入休眠状态。
# 小团手里总共有电量b，他准备在0时刻的时候将所有他想启动的机器人全部一起启动，
# 过后不再启动新的机器人，并且启动的机器人的电量之和不大于b。
# 他想知道在最佳选择的情况下，最快多久能收完衣服。若无论如何怎样都收不完衣服,输出-1.
# !n<=1000 pi<=100 ti,b<=1e5

# !二分答案+dp O(n^2log1e8)

n, b = map(int, input().split())
powers = list(map(int, input().split()))
times = list(map(int, input().split()))

if powers[0] > b:
    print(-1)
    exit(0)
