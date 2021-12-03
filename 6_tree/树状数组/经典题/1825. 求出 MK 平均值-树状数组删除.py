from collections import deque

# MK 平均值 按照如下步骤计算：

# 如果数据流中的整数少于 m 个，MK 平均值 为 -1 ，否则将数据流中最后 m 个元素拷贝到一个独立的容器中。
# 从这个容器中删除最小的 k 个数和最大的 k 个数。
# 计算剩余元素的平均值，并 向下取整到最近的整数 。


class MKAverage:
    # 用一个空的数据流和两个整数 m 和 k 初始化 MKAverage 对象。
    def __init__(self, m: int, k: int):
        self.m = m
        self.k = k
        self.data = deque()
        self.sum_tree = FenwickTree(100000)
        self.count_tree = FenwickTree(100000)

    # 往数据流中插入一个新的整数 num ,1 <= num <= 105
    # 树状数组维护
    def addElement(self, num: int) -> None:
        self.data.append(num)
        self.sum_tree.add(num, num)
        self.count_tree.add(num, 1)
        if len(self.data) > self.m:
            popnum = self.data.popleft()
            # 树状数组删除
            self.sum_tree.add(popnum, -popnum)
            self.count_tree.add(popnum, -1)

    # 对当前的数据流计算并返回 MK 平均数 ，结果需 向下取整到最近的整数 。
    # 前缀和相减
    def calculateMKAverage(self) -> int:
        if len(self.data) < self.m:
            return -1
        left = self.__find_k(self.k)  # 第k小的数
        right = self.__find_k(self.m - self.k)  # 最k大的数
        sum1 = self.sum_tree.query(left)
        sum2 = self.sum_tree.query(right)
        smaller_count1 = self.count_tree.query(left)
        smaller_count2 = self.count_tree.query(right)

        # 减的原因是因为有重复的元素都等于left/right 多算了 要减去这些没有被删除的
        sum1 -= (smaller_count1 - self.k) * left
        sum2 -= (smaller_count2 - (self.m - self.k)) * right
        return (sum2 - sum1) // (self.m - self.k * 2)

    # 寻找第k小的数
    def __find_k(self, k):
        left, right = 0, 100000
        while left <= right:
            mid = (left + right) >> 1
            if self.count_tree.query(mid) < k:
                left = mid + 1
            else:
                right = mid - 1
        return left


class FenwickTree:
    def __init__(self, n):
        self.size = n
        self.tree = [0 for _ in range(n + 1)]  # 索引从1开始，索引0不记录内容

    @staticmethod
    def __lowbit(index):
        return index & -index

    def add(self, index, delta):
        while index <= self.size:
            self.tree[index] += delta
            index += self.__lowbit(index)

    def query(self, index):
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= self.__lowbit(index)
        return res
