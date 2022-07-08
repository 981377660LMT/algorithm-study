from typing import List, Tuple
from collections import defaultdict, Counter, deque

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def minTransfers(self, distributions: List[List[int]]) -> int:
        # 1 <= distributions.length <= 8
        # distributions[i].length == 3
        # 若要使得每一个门店最终借出和借入的商品数量相同，请问至少还需要进行多少次商品调配。
        borrow = defaultdict(int)
        for u, v, w in distributions:
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
