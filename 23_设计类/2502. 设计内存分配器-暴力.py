# 6259. 设计内存分配器
# 给你一个整数 n ，表示下标从 0 开始的内存数组的大小。所有内存单元开始都是空闲的。
# 请你设计一个具备以下功能的内存分配器：
# 1. 分配 一块大小为 size 的连续空闲内存单元并赋 id mID 。
# 2. 释放 给定 id mID 对应的所有内存单元。

# 内存分配/CPU/日程安排


# !1. 优化:线段树O(nlogn)
# https://leetcode.cn/problems/design-memory-allocator/solutions/2024437/bing-mei-you-geng-kuai-de-wen-ding-lgnfe-ze8b/


class Allocator:
    def __init__(self, n: int):
        self.visited = [-1] * n

    def allocate(self, size: int, mID: int) -> int:
        """
        找出大小为 size 个连续空闲内存单元且位于  最左侧 的块，分配并赋 id mID 。
        返回块的第一个下标。如果不存在这样的块，返回 -1 。
        """
        dp = 0
        for i, id in enumerate(self.visited):
            if id == -1:
                dp += 1
                if dp == size:
                    self.visited[i - size + 1 : i + 1] = [mID] * size
                    return i - size + 1
            else:
                dp = 0
        return -1

    def free(self, mID: int) -> int:
        """
        释放 id mID 对应的所有内存单元。返回释放的内存单元数目。
        """
        res = 0
        for i in range(len(self.visited)):
            if self.visited[i] == mID:
                self.visited[i] = -1
                res += 1
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
