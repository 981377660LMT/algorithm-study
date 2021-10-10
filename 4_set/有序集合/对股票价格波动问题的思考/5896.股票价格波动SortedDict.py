# // 请你设计一个算法，实现：

# // 更新 股票在某一时间戳的股票价格，如果有之前同一时间戳的价格，这一操作将 更正 之前的错误价格。
# // 找到当前记录里 最新股票价格 。最新股票价格 定义为时间戳最晚的股票价格。
# // 找到当前记录里股票的 最高价格 。
# // 找到当前记录里股票的 最低价格 。

from sortedcontainers import SortedDict
from collections import defaultdict


class StockPrice:
    def __init__(self):
        self.cur_time = 0
        self.stock_record = defaultdict(int)
        self.prices_record = SortedDict()

    def update(self, timestamp: int, price: int) -> None:
        if timestamp in self.stock_record:
            oldprice = self.stock_record[timestamp]
            self.prices_record[oldprice] -= 1
            if self.prices_record[oldprice] == 0:
                self.prices_record.pop(oldprice)
        self.prices_record.setdefault(price, self.prices_record.get(price, 0) + 1)
        self.cur_time = max(self.cur_time, timestamp)
        self.stock_record[timestamp] = price

    def current(self) -> int:
        return self.stock_record[self.cur_time]

    def maximum(self) -> int:
        return self.prices_record.peekitem()[0]

    def minimum(self) -> int:
        return self.prices_record.peekitem(0)[0]


if __name__ == '__main__':
    stp = StockPrice()
    stp.update(1, 10)
    stp.update(2, 5)
    print(stp.current())
    print(stp.maximum())
    print(stp.update(1, 3))
    print(stp.maximum())
    print(stp.update(4, 2))
    print(stp.minimum())
