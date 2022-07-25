from collections import defaultdict
from sortedcontainers import SortedList


class NumberContainers:
    def __init__(self):
        # !发现还是 ...To... 这种命名看的最习惯
        self.indexToValue = dict()
        self.indexMap = defaultdict(SortedList)

    def change(self, index: int, number: int) -> None:
        """在下标 index 处填入 number 。如果该下标 index 处已经有数字了，那么用 number 替换该数字。"""
        if index in self.indexToValue:
            toRemove = self.indexToValue[index]
            self.indexMap[toRemove].remove(index)
        self.indexToValue[index] = number
        self.indexMap[number].add(index)

    def find(self, number: int) -> int:
        """返回给定数字 number 在系统中的最小下标。如果系统中没有 number ，那么返回 -1 。"""
        if not self.indexMap[number]:
            return -1
        return self.indexMap[number][0]


# Your NumberContainers object will be instantiated and called as such:
# obj = NumberContainers()
# obj.change(index,number)
# param_2 = obj.find(number)

# NumberContainers nc = new NumberContainers();
# nc.find(10); // 没有数字 10 ，所以返回 -1 。
# nc.change(2, 10); // 容器中下标为 2 处填入数字 10 。
# nc.change(1, 10); // 容器中下标为 1 处填入数字 10 。
# nc.change(3, 10); // 容器中下标为 3 处填入数字 10 。
# nc.change(5, 10); // 容器中下标为 5 处填入数字 10 。
# nc.find(10); // 数字 10 所在的下标为 1 ，2 ，3 和 5 。因为最小下标为 1 ，所以返回 1 。
# nc.change(1, 20); // 容器中下标为 1 处填入数字 20 。注意，下标 1 处之前为 10 ，现在被替换为 20 。
# nc.find(10); // 数字 10 所在下标为 2 ，3 和 5 。最小下标为 2 ，所以返回 2 。
nc = NumberContainers()
print(nc.find(10))
nc.change(2, 10)
nc.change(1, 10)
nc.change(3, 10)
nc.change(5, 10)
print(nc.find(10))
nc.change(1, 20)
print(nc.find(10))
# [[],[1,10],[10],[1,20],[10],[20],[30]]
# ["NumberContainers","find","change","change","change","change","find","change","find"]
# [[],[10],[2,10],[1,10],[3,10],[5,10],[10],[1,20],[10]]
# ["NumberContainers","change","find","change","find","find","find"]
# [[],[1,10],[10],[1,20],[10],[20],[30]]
