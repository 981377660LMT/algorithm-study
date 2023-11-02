# https://www.luogu.com.cn/problem/CF3B
# !所有物品只有两种重量(1和2)的背包问题。求不超过背包容量的最大价值以及选择的物品。
# 当其中一种重量的物品数量确定时，另一种重量的物品数量上界可以被确定。
# 对于同一种重量的所有物品，当知道选择个数为x时，可以 贪心 地选择价值最大的x个。这个最优价值经排序预处理后可以得知。
# 因此支持我们枚举重量为a的物品选多少件，对所有情况的最优策略取最大值。


from collections import defaultdict
from itertools import accumulate
from typing import List, Tuple


def lorry(values: List[int], weights: List[int], maxCapacity: int) -> Tuple[int, List[int]]:
    groups = defaultdict(list)  # 按照重量分组
    for i, (value, weight) in enumerate(zip(values, weights)):
        groups[weight].append((value, i))
    for g in groups.values():
        g.sort(key=lambda x: -x[0])
    preSum = defaultdict(lambda: [0])
    for weight, group in groups.items():
        preSum[weight] = [0] + list(accumulate(v for v, _ in group))

    bestCount1 = 0
    maxRes = -1
    w1, w2 = 1, 2  # !两种重量
    n1, n2 = len(groups[w1]), len(groups[w2])
    for count1 in range(min(maxCapacity // w1, n1) + 1):
        count2 = min((maxCapacity - count1 * w1) // w2, n2)
        cand = preSum[w1][count1] + preSum[w2][count2]
        if cand > maxRes:
            maxRes = cand
            bestCount1 = count1

    bestCount2 = min((maxCapacity - bestCount1 * w1) // w2, n2)
    return maxRes, [i for _, i in groups[w1][:bestCount1]] + [i for _, i in groups[w2][:bestCount2]]


if __name__ == "__main__":
    import sys

    input = sys.stdin.readline

    n, maxCapacity = map(int, input().split())
    values = []
    weights = []
    for i in range(n):
        w, v = map(int, input().split())
        values.append(v)
        weights.append(w)

    res, select = lorry(values, weights, maxCapacity)
    print(res)
    for i in select:
        print(i + 1, end=" ")

    # import random

    # for _ in range(1000):
    #     n = random.randint(1, 10)
    #     values = [random.randint(1, 100) for _ in range(n)]
    #     weights = [random.randint(1, 2) for _ in range(n)]
    #     maxCapacity = random.randint(1, 1000)
    #     lorry(values, weights, maxCapacity)
