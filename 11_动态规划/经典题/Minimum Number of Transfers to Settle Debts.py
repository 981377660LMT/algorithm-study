from typing import List
from collections import defaultdict

# 可以用一个三元组 (x, y, z) 表示一次交易，表示 x 借给 y 共计 z 美元
# 给定一群人之间的交易信息列表，计算能够还清所有债务的最小次数。


# 上帝视角两组人，一组欠别人钱，一组别人欠他。两组中相等的数值可以直接匹配。剩下的搜索。
# https://leetcode-cn.com/problems/optimal-account-balancing/comments/336080

# n ≤ 13
def minTransfers(transactions: List[List[int]]) -> int:
    borrow = defaultdict(int)
    for u, v, w in transactions:
        borrow[u] += w
        borrow[v] -= w

    accounts = sorted(list(borrow.values()), reverse=True)

    res = int(1e20)

    def bt(cur: int, times: int):
        nonlocal res
        if times >= res:
            return

        # 账号为0不考虑
        while cur < len(accounts) and accounts[cur] == 0:
            cur += 1

        if cur == len(accounts):
            res = min(res, times)
            return

        for next in range(cur + 1, len(accounts)):
            if accounts[cur] * accounts[next] < 0:
                # 这里加表示cur 给 next 钱 用 [5,-10,5]来想
                accounts[next] += accounts[cur]
                bt(cur + 1, times + 1)
                accounts[next] -= accounts[cur]

    bt(0, 0)
    return res


# 还清债款的最小交易次数
# n ≤ 13
class Solution:
    def solve(self, transfers):
        """Return the minimum number of person-to-person transfers that are required so that all debts are paid."""
        return minTransfers(transfers)


print(Solution().solve(transfers=[[0, 1, 50], [1, 2, 50]]))
# Person 0 gave person 1 $50 and person 1 gave person 2 $50.
# So person 2 can directly give $50 to person 0.
