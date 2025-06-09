# 3573. 买卖股票的最佳时机 V
# https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-v/description/
# 你最多可以进行 k 笔交易，每笔交易可以是以下任一类型：
# 普通交易：在第 i 天买入，然后在之后的第 j 天卖出，其中 i < j。你的利润是 prices[j] - prices[i]。
# 做空交易：在第 i 天卖出，然后在之后的第 j 天买回，其中 i < j。你的利润是 prices[i] - prices[j]。
# 注意：你必须在开始下一笔交易之前完成当前交易。此外，你不能在已经进行买入或卖出操作的同一天再次进行买入或卖出操作。
# !通过进行 最多 k 笔交易，返回你可以获得的最大总利润。
#
# 2 <= prices.length <= 1e3
#
# !做空(Short selling transaction)：预期某个资产（如股票）价格会下跌，通过借入并卖出该资产，等价格下跌后再买回来归还，从而赚取差价。
# 简单来说，就是先卖后买，赚取价格下跌的利润。
#
# 金融和股票交易的专业术语：
#
# - **开多仓**：买入股票，持有多头仓位，预期股价上涨后卖出获利。
#   “开”表示新建仓位，“多仓”指持有实物（股票）。
#
# - **开空仓**：卖出借来的股票，持有空头仓位，预期股价下跌后买回归还获利。
#   “开”表示新建仓位，“空仓”指持有负的股票（即借来卖出）。
#
# - **平多仓**：卖出手中的股票，结束多头仓位。
#   “平”表示平掉已有仓位，“多仓”指之前持有的股票。
#
# - **平空仓**：买入股票归还，结束空头仓位。
#   “平”表示平掉已有仓位，“空仓”指之前借来卖出的股票。
#
# 简言之，“开”是新建仓位，“平”是结束仓位；“多仓”是看涨持有，“空仓”是看跌做空。

from typing import List
from functools import lru_cache


INF = int(1e20)


class Solution:
    def maximumProfit(self, prices: List[int], k: int) -> int:
        @lru_cache(None)
        def dfs(i: int, j: int, endState: int) -> int:
            """ensState: 0-未持有股票(空仓), 1-持有股票(持有多仓), 2-做空股票(持有空仓)."""
            if j > k:
                return -INF
            if i == len(prices):
                return 0 if endState == 0 else -INF

            p = prices[i]
            res = dfs(i + 1, j, endState)  # 不操作
            if endState == 0:
                # 开多仓
                res = max(res, dfs(i + 1, j, 1) - p)
                # 开空仓
                res = max(res, dfs(i + 1, j, 2) + p)
            elif endState == 1:
                # 平多仓
                res = max(res, dfs(i + 1, j + 1, 0) + p)
            elif endState == 2:
                # 平空仓
                res = max(res, dfs(i + 1, j + 1, 0) - p)

            return res

        res = dfs(0, 0, 0)
        dfs.cache_clear()
        return res
