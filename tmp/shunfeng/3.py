# 顺丰里面每天物流信息非常多，为了方便管理物流信息，顺丰的后台有一个物流管理系统，现在将物流系统接口函数简化如下，请你完成如下函数来实现简单的物流管理系统:

# List<ExpressDelivery> query(int pageNo, int pageSize): 按照订单编号逆序查询，每页pageSize条数据，返回逆序查询结果的第pageNo页的物流数据

# ExpressDelivery searchByOrderId(int orderId): 根据订单id搜索物流信息，该方法要求O(1)或者O(logN)查询。如果没有该orderId则返回空

# List<ExpressDelivery> searchByStart(String start): 根据物流起点查对应的订单， 按照订单编号逆序返回。如果没有则返回空数组


# 	int add(String name,String start, String end): 增加订单，返回订单编号。订单编号从1开始自增


# 	你只需要完成上面这几个函数即可，其余的输入输出处理已经写好


# 	保证需要输出对象不超过200000个


from collections import defaultdict, deque
import sys

from typing import List

input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


class ExpressDelivery:
    __slots__ = ("orderId", "name", "start", "end")

    def __init__(self, orderId=-1, name="", start="", end=""):
        self.orderId = orderId  # 订单编号
        self.name = name  # 司机姓名
        self.start = start  # 物流起点
        self.end = end  # 物流终点

    def getOrderId(self):
        return self.orderId

    def getName(self):
        return self.name

    def getStart(self):
        return self.start

    def getEnd(self):
        return self.end


class Main:
    def __init__(self):
        self._id = 1
        self._idToOrder = defaultdict(ExpressDelivery)
        self._startToOrders = defaultdict(list)
        self._hitory = []

    def query(self, pageNo, pageSize) -> List["ExpressDelivery"]:
        # 按照订单编号逆序查询，每页pageSize条数据，返回逆序查询结果的第pageNo页的物流数据
        # 逆序查询第1页=>
        start = pageSize * (pageNo - 1)
        end = pageSize * pageNo
        return self._hitory[start:end][::-1]

    def searchByOrderId(self, orderId) -> "ExpressDelivery":
        if orderId in self._idToOrder:
            return self._idToOrder[orderId]
        return None

    def searchByStart(self, start) -> List["ExpressDelivery"]:
        return self._startToOrders[start][::-1]

    def add(self, name, start, end) -> int:
        id_ = self._id
        newOrder = ExpressDelivery(id_, name, start, end)
        self._idToOrder[id_] = newOrder
        self._startToOrders[start].append(newOrder)
        self._hitory.append(newOrder)
        self._id += 1
        return id_


# 说明： 第一个searchByStartbeijing可以得到北京为出发地的信息，有2个，根据id逆向输出:2_lisi_beijing_hangzhou1_zhangsan_beijing_shanghai第二个searchByOrderId2查找订单号为2的订单:2_lisi_beijing_hangzhou最后query1 20得到4_fu_beijing_hangzhou3_wangsu_hangzhou_shanghai2_lisi_beijing_hangzhou1_zhangsan_beijing_shanghai
# 7
# add
# zhangsan beijing shanghai
# add
# lisi beijing hangzhou
# add
# wangsu hangzhou shanghai
# searchByStart
# beijing
# searchByOrderId
# 2
# add
# fu beijing hangzhou
# query
# 1 20
if __name__ == "__main__":
    main = Main()
    printExpress = []
    opNum = int(input())
    while opNum > 0:
        op = input().strip()
        data = input().strip()
        if op == "add":
            name, start, end = data.split()
            main.add(name, start, end)
        elif op == "searchByStart":
            deliveries = main.searchByStart(data)
            if deliveries:
                printExpress.extend(deliveries)
        elif op == "searchByOrderId":
            orderId = int(data)
            delivery = main.searchByOrderId(orderId)
            if delivery:
                printExpress.append(delivery)
        elif op == "query":
            pageNo, pageSize = map(int, data.split())
            deliveries = main.query(pageNo, pageSize)
            if deliveries:
                printExpress.extend(deliveries)
        opNum -= 1
    for delivery in printExpress:
        deliveryStr = f"{delivery.getOrderId()}_{delivery.getName()}_{delivery.getStart()}_{delivery.getEnd()}"
        print(deliveryStr)
