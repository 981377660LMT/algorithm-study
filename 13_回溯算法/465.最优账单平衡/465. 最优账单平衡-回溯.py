from typing import List
from collections import defaultdict

# 可以用一个三元组 (x, y, z) 表示一次交易，表示 x 借给 y 共计 z 美元
# 给定一群人之间的交易信息列表，计算能够还清所有债务的最小次数。


# 上帝视角两组人，一组欠别人钱，一组别人欠他。两组中相等的数值可以直接匹配。剩下的搜索。
# https://leetcode-cn.com/problems/optimal-account-balancing/comments/336080
class Solution:
    def minTransfers(self, transactions: List[List[int]]) -> int:
        def bt(index: int, curSum: int) -> None:
            nonlocal res
            if curSum >= res:
                return

            while index < len(deg) and deg[index] == 0:  # 去掉相等的数值
                index += 1

            if index == len(deg):
                res = min(res, curSum)
                return

            for next in range(index + 1, len(deg)):  # 搜索这个人和谁交易 要把钱全部交易完
                if deg[index] * deg[next] < 0:
                    deg[next] += deg[index]
                    bt(index + 1, curSum + 1)
                    deg[next] -= deg[index]

        deg = defaultdict(int)
        for u, v, w in transactions:
            deg[u] += w
            deg[v] -= w
        deg = sorted([v for v in deg.values() if v != 0], reverse=True)
        res = int(1e20)
        bt(0, 0)
        return res


print(Solution().minTransfers([[0, 1, 10], [2, 0, 5]]))
# 输出：
# 2

# 解释：
# 人 #0 给人 #1 共计 10 美元。
# 人 #2 给人 #0 共计 5 美元。

# 需要两次交易。一种方式是人 #1 分别给人 #0 和人 #2 各 5 美元。
