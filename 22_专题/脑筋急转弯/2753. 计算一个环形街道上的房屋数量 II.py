# 2753. 计算一个环形街道上的房屋数量 II

# 每个房屋的门初始时可以是开着的也可以是关着的（至少有一个房屋的门是开着的）。
# !房屋数量不超过k
# 1 <= n <= k <= 1e5
# 求房屋数量

# 1.不断向右找到一扇开着的门。
# 2.跳过这扇门，继续向右。每向右一步，计数器 i 就加一。循环至多 k 次。
# 3.如果遇到开着的门，就把门关上。
# 4.答案为最后一次遇到的开着的门的 i。


class IStreet:
    def closeDoor(self):
        pass

    def isDoorOpen(self):
        pass

    def moveRight(self):
        pass


class Solution:
    def houseCount(self, s: "IStreet", k: int) -> int:
        while not s.isDoorOpen():
            s.moveRight()
        res = 0
        for i in range(k):
            s.moveRight()
            if s.isDoorOpen():
                s.closeDoor()
                res = i
        return res + 1
