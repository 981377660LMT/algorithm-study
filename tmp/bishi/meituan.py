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
# 他有n件衣服从左到右成一行排列，所在位置分别为1~n,
# 在每个位置上已经有一个就绪的机器人可以帮忙收衣服，但第i个位置上的机器人需要pi的电量来启动。
# 然后这个机器人会用ti的时间收衣服，当它收完当前衣服后，
# 会尝试去收紧邻的右边的一件衣服(如果存在的话)，即i+1处的衣服，
# 如果i+1处的衣服已经被其他机器人收了或者其他机器人正在收，这个机器人就会进入休眠状态，不再收衣服。
# 不过如果机器人没有休眠，它会同样以ti时间来收这件1+1处的衣服(注意，不是t+1的时间，收衣服的时间为每个机器人固有属性)，
# 然后它会做同样的检测来看能否继续收i+2处的衣服，一直直到它进入休眠状态或者右边没有衣服可以收了。
# 形象地来说，机器人会一直尝试往右边收衣服，收k件的话就耗费k*ti;的时间,
# 但是当它遇见其他机器人工作的痕迹，就会认为后面的事情它不用管了,开始摸鱼，进入休眠状态。
# !小团手里总共有电量b，他准备在0时刻的时候将所有他想启动的机器人全部一起(并行)启动，
# 过后不再启动新的机器人，并且启动的机器人的电量之和不大于b。
# 他想知道在最佳选择的情况下，最快多久能收完衣服。若无论如何怎样都收不完衣服,输出-1.
# !n<=1000 pi<=100 ti,b<=1e5

# !二分答案+dp O(n^2log1e8)
# !为什么想到二分:因为需要固定一个变量(题目有两个变量:时间和电量)
# 首先二分时间mid，然后dp, dp[i]表示在不超过mid的时间内收完前i个所需要的最少启动电量
# 求解dp[i]时，往前枚举j，如果times[j+1]*(i -j) <= mid，
# 则 dp[i] = min(dp[i], dp[j]+ powers[j+1])
# 然后看最后的dp[n]<= b


from typing import List

INF = int(4e18)


def assignRobots(n: int, b: int, powers: List[int], times: List[int]) -> int:
    """美团笔试t4

    Args:
        n (int): n件衣服从左到右成一行排列
        b (int): 启动的机器人的电量之和不大于b
        powers (List[int]): 启动每个机器人的花费
        times (List[int]): 每个机器人收一件衣服的时间

    Returns:
        int: 在最佳选择的情况下，最快多久能收完衣服;若无论如何怎样都收不完衣服,输出-1
    """

    def check(mid: int) -> bool:
        """能否在不超过mid的时间内收完"""
        dp = [INF] * (n + 1)  # dp[i]表示在不超过mid的时间内收完前i个所需要的最少启动电量
        dp[0] = 0
        for i in range(1, n + 1):
            for j in range(i):  # 可以线段树优化 但是要反过来用i去更新区段i+1到 i + mid / times[j] 的最小值
                if times[j] * (i - j) <= mid:
                    dp[i] = min(dp[i], dp[j] + powers[j])
        return dp[-1] <= b

    left, right = 0, int(1e8 + 10)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            right = mid - 1
        else:
            left = mid + 1

    return left if check(left) else -1


print(assignRobots(n=3, b=2, powers=[1, 1, 3], times=[10, 2, 3]))  # 10
print(assignRobots(n=4, b=2, powers=[1, 10, 1, 10], times=[10, 2, 5, 1]))  # 20
