from typing import List


MONEY = [20, 50, 100, 200, 500]


class ATM:
    def __init__(self):
        # 根本不用counter，不存钱而是存index
        self.counter = [0, 0, 0, 0, 0]

    def deposit(self, banknotesCount: List[int]) -> None:
        # 分别存入 $20 ，$50，$100，$200 和 $500 钞票的数目。
        for i in range(5):
            self.counter[i] += banknotesCount[i]

    def withdraw(self, amount: int) -> List[int]:
        # 如果无法取出指定数额的钱，请返回 [-1] （这种情况下 不 取出任何钞票）。
        res = [0, 0, 0, 0, 0]
        for i in range(4, -1, -1):
            count = self.counter[i]
            div, _ = divmod(amount, MONEY[i])
            ok = min(count, div)
            res[i] += ok
            amount -= ok * MONEY[i]

        if amount != 0:
            return [-1]
        else:
            for i in range(5):
                self.counter[i] -= res[i]
            return res


atm = ATM()
atm.deposit([0, 0, 1, 2, 1])
print(atm.withdraw(600))
