# 银联-4. 设计自动售货机
# !1.注意遍历数组时不能修改遍历中的数组 (需要遍历数组的拷贝)
# !2.注意不能直接修改SortedList某个位置上的元素 需要先删除再插入
# 时间复杂度:O(nq)

from math import ceil
from collections import defaultdict
from sortedcontainers import SortedList


class VendingMachine:
    def __init__(self):
        # 同一名顾客每成功购买一次，下次购买便可多享受 1% 的折扣（折后价向上取整），
        # 最低折扣为 70%
        self.discount = defaultdict(lambda: 100)  # 用户折扣
        self.goods = defaultdict(SortedList)  # 商品列表

    def addItem(self, time: int, number: int, item: str, price: int, duration: int) -> None:
        """在 time 时刻向售货机中增加 number 个名称为 item 的商品，价格为 price,保质期为 duration"""
        self.goods[item].add((price, time + duration, number))

    def sell(self, time: int, customer: str, item: str, number: int) -> int:
        """在 time 时刻，名称为 customer 的顾客前来购买了 number 个名称为 item 的商品，返回总费用

        当且仅当售货机中存在足够数量的未过期商品方可成功购买，并返回支付的总费用，
        否则一件商品也不会售出，并返回 -1
        对于价格不同的同种商品，优先售出价格最低的商品；
        如果有价格相同的同种商品，优先出售距离过期时间最近的商品；
        """
        items = self.goods[item]

        remain = 0
        for price, expire, count in items.copy():
            if expire < time:
                items.remove((price, expire, count))
            else:
                remain += count

        if remain < number:
            return -1

        priceSum = 0
        for price, expire, count in items.copy():
            if count >= number:
                next = (price, expire, count - number)
                items.discard((price, expire, count))
                items.add(next)
                priceSum += price * number
                break
            else:
                number -= count
                items.discard((price, expire, count))
                priceSum += price * count

        res = ceil(priceSum * self.discount[customer] / 100)
        self.discount[customer] = max(70, self.discount[customer] - 1)
        return res


if __name__ == "__main__":

    machine = VendingMachine()
    machine.addItem(135, 339, "Inc", 368, 25)
    machine.addItem(146, 377, "Inc", 278, 210)
    machine.addItem(216, 30, "Inc", 492, 48)
    machine.addItem(232, 207, "Inc", 62, 146)
    machine.addItem(356, 215, "Inc", 337, 67)
    print(machine.sell(401, "Timothy", "Inc", 102))
    print(machine.sell(474, "Timothy", "Inc", 480))
