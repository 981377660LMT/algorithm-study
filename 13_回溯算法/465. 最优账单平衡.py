from typing import List
from collections import defaultdict

# 可以用一个三元组 (x, y, z) 表示一次交易，表示 x 借给 y 共计 z 美元
# 给定一群人之间的交易信息列表，计算能够还清所有债务的最小次数。


# 上帝视角两组人，一组欠别人钱，一组别人欠他。两组中相等的数值可以直接匹配。剩下的搜索。
# https://leetcode-cn.com/problems/optimal-account-balancing/comments/336080
class Solution:
    def minTransfers(self, transactions: List[List[int]]) -> int:
        borrow = defaultdict(int)
        for u, v, w in transactions:
            borrow[u] += w
            borrow[v] -= w
        accounts = list(borrow.values())
        print(accounts)
        res = 0x7FFFFFFF

        def dfs(cur: int, times: int):
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
                    dfs(cur + 1, times + 1)
                    accounts[next] -= accounts[cur]

        dfs(0, 0)
        return res


print(Solution().minTransfers([[0, 1, 10], [2, 0, 5]]))
# 输出：
# 2

# 解释：
# 人 #0 给人 #1 共计 10 美元。
# 人 #2 给人 #0 共计 5 美元。

# 需要两次交易。一种方式是人 #1 分别给人 #0 和人 #2 各 5 美元。

