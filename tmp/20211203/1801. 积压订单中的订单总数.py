from typing import List
from heapq import heappop, heappush

# 其中每个 orders[i] = [pricei, amounti, orderTypei] 表示有 amounti 笔类型为 orderTypei 、价格为 pricei 的订单。
# 0 表示这是一批采购订单 buy
# 1 表示这是一批销售订单 sell
# 由 orders[i] 表示的所有订单提交时间均早于 orders[i+1] 表示的所有订单。

# 总结：每次添加订单时都比较两个队列头部，循环相消
class Solution:
    def getNumberOfBacklogOrders(self, orders: List[List[int]]) -> int:
        sell, buy = [], []

        for p, a, t in orders:
            if t == 0:
                heappush(buy, [-p, a])
            else:
                heappush(sell, [p, a])
            while sell and buy and sell[0][0] <= -buy[0][0]:
                k = min(buy[0][1], sell[0][1])
                buy[0][1] -= k
                sell[0][1] -= k
                if buy[0][1] == 0:
                    heappop(buy)
                if sell[0][1] == 0:
                    heappop(sell)

        return sum(a for _, a in (buy + sell)) % (10 ** 9 + 7)


print(Solution().getNumberOfBacklogOrders(orders=[[10, 5, 0], [15, 2, 1], [25, 1, 1], [30, 4, 0]]))


# 描述太长 跳过
# 如果该订单是一笔采购订单 buy ，则可以查看积压订单中价格 最低 的销售订单 sell 。
# 如果该销售订单 sell 的价格 低于或等于 当前采购订单 buy 的价格，
# 则匹配并执行这两笔订单，并将销售订单 sell 从积压订单中删除。
# 否则，采购订单 buy 将会添加到积压订单中。(卖出最赚钱的)

# 反之亦然，如果该订单是一笔销售订单 sell ，
# 则可以查看积压订单中价格 最高 的采购订单 buy 。
# 如果该采购订单 buy 的价格 高于或等于 当前销售订单 sell 的价格，
# 则匹配并执行这两笔订单，并将采购订单 buy 从积压订单中删除。
# 否则，销售订单 sell 将会添加到积压订单中。(买入贵的)

# 总结:sell始终要低于buy的价格

# 输入所有订单后，返回积压订单中的 订单总数 。由于数字可能很大，所以需要返回对 109 + 7 取余的结果。
