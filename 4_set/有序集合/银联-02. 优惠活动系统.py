from collections import defaultdict
from sortedcontainers import SortedList


# 这种题目都是sortedList存元组，另外用一个哈希表存元组id=>元组
# 此外还要关注题目的主体是谁，比如用户，就给用户开一个 defaultdict，商店，就给商店开一个 defaultdict
class DiscountSystem:
    def __init__(self):
        self.activity = SortedList()
        self.actById = defaultdict(tuple)
        # 每个用户活动的参加次数
        self.userInfo = defaultdict(lambda: defaultdict(int))

    def addActivity(
        self, actId: int, priceLimit: int, discount: int, number: int, userLimit: int
    ) -> None:
        """
        单笔消费的原价不小于 priceLimit 时，可享受 discount 的减免
        每个用户最多参与该优惠活动 userLimit 次
        该优惠活动共有 number 数量的参加名额
        """
        cur = (-discount, actId, priceLimit, userLimit, number)
        self.activity.add(cur)
        self.actById[actId] = cur

    def removeActivity(self, actId: int) -> None:
        """表示结束编号为 actId 的优惠活动"""
        act = self.actById[actId]
        self.activity.discard(act)

    def consume(self, userId: int, cost: int) -> int:
        """
        表示用户 userId 产生了一笔原价为 cost 的消费。请返回用户的实际支付金额。
        单次消费最多可参加一份优惠活动
        若可享受优惠减免，则 「支付金额 = 原价 - 优惠减免」
        若同时满足多个优惠活动时，则优先参加优惠减免最大的活动
        若有多个优惠减免最大的活动，优先参加 actId 最小的活动
        """

        for i in range(len(self.activity)):
            raw = self.activity[i]

            discount, actId, priceLimit, userLimit, number = raw
            discount *= -1

            if cost < priceLimit:
                continue
            if number <= 0:
                continue

            used = self.userInfo[userId][actId]
            if used >= userLimit:
                continue

            self.userInfo[userId][actId] += 1
            number -= 1

            self.activity.discard(raw)
            newItem = (-discount, actId, priceLimit, userLimit, number)
            self.activity.add(newItem)
            self.actById[actId] = newItem

            return cost - discount

        return cost


# [null,null,11,null,19]
# [null, null, null, 7, null, 11, 11, 10, 21]
