from typing import List
from heapq import heappush, heappushpop


class MedianFinder:
    def __init__(self):
        self.small = []  # 大顶堆
        self.large = []  # 小顶堆
        self.leftSum = 0
        self.rightSum = 0

    def add(self, num: int) -> None:
        if len(self.small) == len(self.large):
            leftMax = -heappushpop(self.small, -num)
            heappush(self.large, leftMax)
            self.leftSum += num - leftMax
            self.rightSum += leftMax
        elif len(self.small) < len(self.large):
            rightMin = heappushpop(self.large, num)
            heappush(self.small, -rightMin)
            self.leftSum += rightMin
            self.rightSum += num - rightMin

    def getMedian(self) -> float:
        if len(self.small) == len(self.large):
            return (self.large[0] - self.small[0]) / 2
        else:
            return self.large[0]

    # 所有数到中位数的距离(绝对值)之和
    def getDiffSum(self) -> float:
        median = self.getMedian()
        return self.rightSum - len(self.large) * median + len(self.small) * median - self.leftSum


# 首先将数据进行等效转换以方便计算：g[i] = nums[i] - i。
# 转换完之后题目变为：给你一个数列g，求把g的所有元素变成同一个值的最小代价，那么显然这个值应该是g的中位数。
# 使用双堆维护中位数以及中位数两边的数据累加和，便于快速累加差值。

# 主办方请小扣回答出一个长度为 N 的数组，
# 第 i 个元素(0 <= i < N)表示将 0~i 号计数器 初始 所示数字操作成满足所有条件
# nums[a]+1 == nums[a+1],(0 <= a < i) 的最小操作数。

# 1703. 得到连续 K 个 1 的最少相邻交换次数
# 假设变化前的值为[o0,o1,...,ok-1]，变化后为[x,x+1,...,x+k-1]
# 那么要求的就是|o0-x|+|o1-(x+1)|+...+|ok-1-(x+k-1)|最小值，变形得
# |o0-x|+|(o1-1)-x|+|(ok-1-(k-1))-x| 即这个x就是他们的中位数mid
class Solution:
    def numsGame(self, nums: List[int]) -> List[int]:
        MOD = int(1e9 + 7)
        res = []
        medianFinder = MedianFinder()
        for i in range(len(nums)):
            # 平移
            normalized = nums[i] - i
            medianFinder.add(normalized)
            res.append(int(medianFinder.getDiffSum() % MOD))
        # print(medianFinder.__dict__)
        return res


print(Solution().numsGame([1, 1, 1, 2, 3, 4]))
# 输出：[0,1,2,3,3,3]

# 解释：
# i = 0，无需操作；
# i = 1，将 [1,1] 操作成 [1,2] 或 [0,1] 最少 1 次操作；
# i = 2，将 [1,1,1] 操作成 [1,2,3] 或 [0,1,2]，最少 2 次操作；
# i = 3，将 [1,1,1,2] 操作成 [1,2,3,4] 或 [0,1,2,3]，最少 3 次操作；
# i = 4，将 [1,1,1,2,3] 操作成 [-1,0,1,2,3]，最少 3 次操作；
# i = 5，将 [1,1,1,2,3,4] 操作成 [-1,0,1,2,3,4]，最少 3 次操作；
# 返回 [0,1,2,3,3,3]。
