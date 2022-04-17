from collections import Counter
from typing import List, Optional, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class ATM:
    def __init__(self):
        self.counter = Counter()

    def deposit(self, banknotesCount: List[int]) -> None:
        # 分别存入 $20 ，$50，$100，$200 和 $500 钞票的数目。
        m20, m50, m100, m200, m500 = banknotesCount
        self.counter[20] += m20
        self.counter[50] += m50
        self.counter[100] += m100
        self.counter[200] += m200
        self.counter[500] += m500

    def withdraw(self, amount: int) -> List[int]:
        # 如果无法取出指定数额的钱，请返回 [-1] （这种情况下 不 取出任何钞票）。
        res = [0, 0, 0, 0, 0]
        for i, money in enumerate([500, 200, 100, 50, 20], start=1):
            count = self.counter[money]
            div, _ = divmod(amount, money)
            ok = min(count, div)
            res[-i] += ok
            amount -= ok * money

        if amount != 0:
            return [-1]
        else:
            for i, money in enumerate([500, 200, 100, 50, 20], start=1):
                self.counter[money] -= res[-i]
            return res


atm = ATM()
atm.deposit([0, 0, 1, 2, 1])
print(atm.withdraw(600))
