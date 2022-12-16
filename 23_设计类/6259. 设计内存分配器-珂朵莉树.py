# 6259. 设计内存分配器
# 给你一个整数 n ，表示下标从 0 开始的内存数组的大小。所有内存单元开始都是空闲的。
# 请你设计一个具备以下功能的内存分配器：
# 1. 分配 一块大小为 size 的连续空闲内存单元并赋 id mID 。
# 2. 释放 给定 id mID 对应的所有内存单元。


# !1. 优化:珂朵莉树(SortedList模拟区间)
# https://leetcode.cn/problems/design-memory-allocator/solution/by-freeyourmind-sc6b/


from sortedcontainers import SortedList


class Allocator:
    def __init__(self, n: int):
        self.sl = SortedList([(-1, 1, -1), (n, 1, -1)])  # SortedList[(start, size, mID)]

    def allocate(self, size: int, mID: int) -> int:
        """
        找出大小为 size 个连续空闲内存单元且位于  最左侧 的块，分配并赋 id mID 。
        返回块的第一个下标。如果不存在这样的块，返回 -1 。
        """
        for (start1, len1, _), (start2, _, _) in zip(self.sl, self.sl[1:]):
            if start2 - (start1 + len1) >= size:
                self.sl.add((start1 + len1, size, mID))
                return start1 + len1
        return -1

    def free(self, mID: int) -> int:
        """
        释放 id mID 对应的所有内存单元。返回释放的内存单元数目。
        """
        toRemove = [i for i, (*_, id) in enumerate(self.sl) if id == mID]
        res = 0
        # !倒着删除不会影响前面需要删除的元素的下标
        for i in toRemove[::-1]:
            item = self.sl.pop(i)
            res += item[1]
        return res


# Your Allocator object will be instantiated and called as such:
# obj = Allocator(n)
# param_1 = obj.allocate(size,mID)
# param_2 = obj.free(mID)
# ["Allocator", "allocate", "allocate", "allocate", "free", "allocate", "allocate", "allocate", "free", "allocate", "free"]
# [[10], [1, 1], [1, 2], [1, 3], [2], [3, 4], [1, 1], [1, 1], [1], [10, 2], [7]]
A = Allocator(10)
print(A.allocate(1, 1))
print(A.allocate(1, 2))
print(A.allocate(1, 3))
print(A.free(2))
print(A.allocate(3, 4))
print(A.allocate(1, 1))
print(A.allocate(1, 1))
print(A.free(1))
print(A.allocate(10, 2))
print(A.free(7))
