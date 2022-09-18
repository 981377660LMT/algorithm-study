# 银联-4. 设计自动售货机
# !一个对象维护两种优先级:价格和过期时间
# !方法是使用两个有序数据结构,一个优先级为购买商品的优先级,一个优先级为过期时间
# (如果需要的话，可以两个有序数据结构共享一个对象的引用，例如堆的懒删除)
# !这样做的好处是，在删除过期的商品时就不用遍历整个商品列表了。
# 时间复杂度:O(qlogq)

from math import ceil
from collections import defaultdict
from typing import Tuple
from sortedcontainers import SortedList


class VendingMachine:
    def __init__(self):
        # 同一名顾客每成功购买一次，下次购买便可多享受 1% 的折扣（折后价向上取整），
        # 最低折扣为 70%
        self.discount = defaultdict(lambda: 100)  # 用户折扣
        self.goods = defaultdict(ItemManager)  # 商品列表

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
        itemManager = self.goods[item]
        itemManager.removeExpired(time)  # 下架过期商品
        priceSum = itemManager.sell(number)
        if priceSum == -1:
            return -1
        res = ceil(priceSum * self.discount[customer] / 100)
        self.discount[customer] = max(70, self.discount[customer] - 1)
        return res


class ItemManager:
    __slots__ = "_data1", "_data2", "_remain"

    def __init__(self) -> None:
        self._data1 = SortedList()  # 优先级为(价格,过期时间) 存储[价格,过期时间,数量]
        self._data2 = SortedList(key=lambda x: x[1])  # 优先级为过期时间 存储[价格,过期时间,数量]
        self._remain = 0

    def add(self, item: Tuple[int, int, int]) -> None:
        """添加商品"""
        self._data1.add(item)
        self._data2.add(item)
        self._remain += item[2]

    def remove(self, item: Tuple[int, int, int]) -> None:
        """下架商品"""
        self._data1.remove(item)
        self._data2.remove(item)
        self._remain -= item[2]

    def sell(self, need: int) -> int:
        """卖出need个商品 返回总价"""
        if need > self._remain:
            return -1

        res = 0
        toRemove, toAdd = [], []
        for item in self._data1:
            price, expire, number = item
            if number <= need:
                res += price * number
                toRemove.append(item)
                need -= number
            else:
                res += price * need
                toRemove.append(item)
                toAdd.append((price, expire, number - need))
                break

        for item in toRemove:
            self.remove(item)

        for item in toAdd:
            self.add(item)

        return res

    def removeExpired(self, time: int) -> None:
        """下架过期时间小于time的商品"""
        while self._data2:
            if self._data2[0][1] >= time:
                break
            self.remove(self._data2[0])  # type: ignore


if __name__ == "__main__":

    machine = VendingMachine()
    machine.addItem(135, 339, "Inc", 368, 25)
    machine.addItem(146, 377, "Inc", 278, 210)
    machine.addItem(216, 30, "Inc", 492, 48)
    machine.addItem(232, 207, "Inc", 62, 146)
    machine.addItem(356, 215, "Inc", 337, 67)
    print(machine.sell(401, "Timothy", "Inc", 102))
    print(machine.sell(474, "Timothy", "Inc", 480))
