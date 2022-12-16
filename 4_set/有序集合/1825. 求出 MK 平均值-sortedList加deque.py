from collections import deque
from sortedcontainers import SortedList

# MK平均值 按照如下步骤计算：

# 1.如果数据流中的整数少于 m 个，MK 平均值 为 -1 ，
# 2.否则将数据流中`最后 m 个元素`拷贝到一个独立的容器中。
# !从这个容器中删除最小的 k 个数和最大的 k 个数。
# 3.计算剩余元素的平均值，并 向下取整到最近的整数 。
# !m<=1e5 k*2<m
# !滑动窗口topK之和
# https://leetcode.cn/problems/finding-mk-average/solution/by-981377660lmt-5hhm/


class MKAverage:
    def __init__(self, m: int, k: int):
        self.windowSum = 0
        self.queue = deque()
        self.sl = SortedList()
        self.m = m
        self.k = k

    def addElement(self, num: int) -> None:
        self.queue.append(num)
        if len(self.queue) == self.m:  # 初始化
            self.sl = SortedList(self.queue)
            self.windowSum = sum(self.sl[self.k : -self.k])

        if len(self.queue) > self.m:
            # 加入后对sum_的影响，num会把sortedList里的元素挤到左边或者右边
            pos = self.sl.bisect_left(num)
            if pos < self.k:
                # 被挤到中间来了
                self.windowSum += self.sl[self.k - 1]  # type: ignore
            elif self.k <= pos <= self.m - self.k:
                self.windowSum += num
            else:
                self.windowSum += self.sl[self.m - self.k]  # type: ignore
            self.sl.add(num)

            popped = self.queue.popleft()
            pos = self.sl.bisect_left(popped)
            if pos < self.k:
                # 左移
                self.windowSum -= self.sl[self.k]  # type: ignore
            elif self.k <= pos <= self.m - self.k:
                self.windowSum -= popped
            else:
                # 右移
                self.windowSum -= self.sl[self.m - self.k]  # type: ignore
            self.sl.remove(popped)

    def calculateMKAverage(self) -> int:
        if len(self.queue) < self.m:
            return -1
        return self.windowSum // (self.m - 2 * self.k)


class TopkSum:
    __slots__ = ("_sl", "_k", "_topKSum")

    def __init__(self, k: int, isMin: bool) -> None:
        self._sl = SortedList() if isMin else SortedList(key=lambda x: -x)
        self._k = k
        self._topKSum = 0

    def add(self, x: int) -> None:
        pos = self._sl.bisect_left(x)
        if pos < self._k:
            self._topKSum += x
            if self._k - 1 < len(self._sl):
                self._topKSum -= self._sl[self._k - 1]  # type: ignore
        self._sl.add(x)

    def remove(self, x: int) -> None:
        pos = self._sl.bisect_left(x)
        if pos < self._k:
            self._topKSum -= x
            if self._k < len(self._sl):
                self._topKSum += self._sl[self._k]  # type: ignore
        self._sl.remove(x)

    def query(self) -> int:
        return self._topKSum


# Your MKAverage object will be instantiated and called as such:
# obj = MKAverage(m, k)
# obj.addElement(num)
# param_2 = obj.calculateMKAverage()
# M, K = 0, 0
# DENO = 1  # 分母
# curSum = 0
# sl = SortedList()
# queue = deque()


# class MKAverage:
#     def __init__(self, m: int, k: int):
#         """用一个空的数据流和两个整数 m 和 k 初始化 MKAverage 对象"""
#         global M, K, DENO, sl, queue, curSum
#         M, K = m, k
#         DENO = M - 2 * K  # 分母，题目保证为正整数
#         curSum = 0
#         sl.clear()
#         queue.clear()

#     def addElement(self, num: int) -> None:
#         """往数据流中插入一个新的整数 num ,1 <= num <= 1e5"""
#         global sl, queue, curSum
#         queue.append(num)

#         if len(queue) == M:
#             # 初始化
#             sl = SortedList(queue)
#             curSum = sum(sl[K:-K])

#         if len(queue) > M:
#             # 加入后对sum_的影响，num会把sortedList里的元素挤到左边或者右边
#             pos = sl.bisect_left(num)
#             if pos < K:
#                 # 被挤到中间来了
#                 curSum += sl[K - 1]
#             elif K <= pos <= M - K:
#                 curSum += num
#             else:
#                 # 被挤到中间来了
#                 curSum += sl[M - K]
#             sl.add(num)

#             # 从deque里出来一个数对sum_的影响
#             popped = queue.popleft()
#             pos = sl.bisect_left(popped)
#             if pos < K:
#                 # 左移
#                 curSum -= sl[K]
#             elif K <= pos <= M - K:
#                 curSum -= popped
#             else:
#                 # 右移
#                 curSum -= sl[M - K]
#             sl.remove(popped)

#     def calculateMKAverage(self) -> int:
#         """对当前的数据流计算并返回 MK 平均数 ，结果需 向下取整到最近的整数 。"""
#         if len(sl) < M:
#             return -1
#         return floor(curSum / DENO)


if __name__ == "__main__":
    MK = MKAverage(3, 1)
    MK.addElement(17612)
    MK.addElement(74607)
    print(MK.calculateMKAverage())
    MK.addElement(8272)
    MK.addElement(33433)
    print(MK.calculateMKAverage())
    MK.addElement(15456)
    MK.addElement(64938)
    print(MK.calculateMKAverage())
