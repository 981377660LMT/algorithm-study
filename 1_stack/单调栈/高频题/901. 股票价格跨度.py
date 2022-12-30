# 股票价格跨度 : 单调栈在线
# 901. 股票价格跨度


class StockSpanner:
    def __init__(self):
        self.stack = []  # (value,day)

    def next(self, price: int) -> int:
        """求股票价格小于或等于今天价格的最大连续日数(包含今天)"""
        count = 1
        while self.stack and self.stack[-1][0] <= price:
            count += self.stack.pop()[1]
        self.stack.append((price, count))
        return count


# Your StockSpanner object will be instantiated and called as such:
# obj = StockSpanner()
# param_1 = obj.next(price)
