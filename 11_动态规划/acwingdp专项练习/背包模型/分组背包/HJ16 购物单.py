# 分组背包
# 如果要买归类为附件的物品，必须先买该附件所属的主件，且每件物品只能购买一次。
# 他希望在花费不超过 N 元的前提下，使自己的满意度达到最大。
# 满意度是指所购买的每件物品的价格与重要度的乘积的总和
# 请你帮助王强计算可获得的最大的满意度。
# money<=32000 表示总钱数
# m<60 表示可购买的物品个数

from collections import defaultdict, deque
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

money, m = map(int, input().split())
goods = defaultdict(deque)
for gid in range(m):
    cost, score, group = map(int, input().split())
    if group == 0:
        goods[gid].appendleft((cost, score))  # 主件
    else:
        goods[group - 1].append((cost, score))  # 附件

keys = sorted(goods)
n = len(keys)


@lru_cache(None)
def dfs(index: int, remain: int) -> int:
    if remain < 0:
        return -INF
    if index == n:
        return 0

    # 不选主件
    res = dfs(index + 1, remain)

    # 选主件
    curGroup = goods[keys[index]]
    for state in range(1 << (len(curGroup) - 1)):
        cur = [0]
        for i in range((len(curGroup) - 1)):
            if state & (1 << i):
                cur.append(i + 1)
        costSum, scoreSum = sum([curGroup[i][0] for i in cur]), sum(
            [curGroup[i][0] * curGroup[i][1] for i in cur]
        )
        res = max(res, scoreSum + dfs(index + 1, remain - costSum))
    return res


print(dfs(0, money))
