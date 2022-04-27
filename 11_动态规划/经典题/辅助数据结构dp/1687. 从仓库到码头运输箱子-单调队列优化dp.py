from itertools import accumulate
from typing import Any, Deque, Iterable, List, Optional
from collections import deque

# 你需要用这一辆车把一些箱子从仓库运送到码头。这辆卡车每次运输有 箱子数目的限制 和 总重量的限制 。
# 箱子需要按照 数组顺序 运输
# 对于在卡车上的箱子，我们需要 按顺序 处理它们

# 请你返回将所有箱子送到相应码头的 最少行程 次数。

# https://leetcode.cn/problems/delivering-boxes-from-storage-to-ports/solution/python-dan-diao-dui-lie-you-hua-dp-by-98-6elh/


class MonoQueue:
    """具有 O(1) 求 `min` 和 `max` API的 deque"""

    def __init__(self, iterable: Optional[Iterable[int]] = None) -> None:
        self.minQueue: Deque[Any] = deque()
        self.maxQueue: Deque[Any] = deque()
        self.rawQueue: Deque[int] = deque()

        if iterable is not None:
            for value in iterable:
                self.append(value)

    @property
    def min(self) -> int:
        if not self.minQueue:
            raise ValueError('monoQueue is empty')
        return self.minQueue[0][0]

    @property
    def max(self) -> int:
        if not self.maxQueue:
            raise ValueError('monoQueue is empty')
        return self.maxQueue[0][0]

    def popleft(self) -> int:
        if not self.rawQueue:
            raise IndexError('popleft from empty queue')

        self.minQueue[0][-1] -= 1
        if self.minQueue[0][-1] == 0:
            self.minQueue.popleft()

        self.maxQueue[0][-1] -= 1
        if self.maxQueue[0][-1] == 0:
            self.maxQueue.popleft()

        return self.rawQueue.popleft()

    def append(self, value: int, *metaInfo: Any) -> None:
        """
        Args:
            value (int): 元素的值
            metaInfo: Any 当前元素附加的元信息，不会添加到原始队列
        """
        count = 1
        while self.minQueue and self.minQueue[-1][0] > value:
            count += self.minQueue.pop()[-1]
        self.minQueue.append([value, *metaInfo, count])

        count = 1
        while self.maxQueue and self.maxQueue[-1][0] < value:
            count += self.maxQueue.pop()[-1]
        self.maxQueue.append([value, *metaInfo, count])

        self.rawQueue.append(value)

    def __len__(self) -> int:
        return len(self.rawQueue)

    def __getitem__(self, index: int) -> int:
        return self.rawQueue[index]


class Solution:
    def boxDelivering(self, boxes: List[List[int]], _: int, maxBoxes: int, maxWeight: int) -> int:
        """问题的关键是哪几个箱子一起运送 dp[i] 表示前i个箱子运送到码头的最少次数
        显然 可以写一个O(n^2) 的dp
        dp[i]=min(dp[j]+cost(j+1,i)+2 for j in range(i)) 其中cost(i+1,j)可以前缀和O(1)求出
        即min(dp[j]+preCost[i]-preCost[j]+2 for j in range(i))

        如何优化这个dp呢 答案就是`删去dp拓扑图中不必要的边`
        对每一个i 不需要从这么多个j转移过来 只需要从最好的一个j转移过来即可
        即用一个数据结构维护之前的 dp[j]-preCost[j] 即可!!!
        这个数据结构可以是一个小根堆，一个平衡树，但是这里只要维护最小值，所以用单调队列(MonoQueue)即可!!!
        """
        n = len(boxes)
        preWeight = [0] + list(accumulate(box[1] for box in boxes))  # 重量前缀和
        preCost = [0]  # 运送次数前缀和 一次搬运前1个到前3个需要`preCost[2] - preCost[0]`次转移
        for (pre, _), (cur, _) in zip(boxes, boxes[1:]):
            preCost.append(preCost[-1] + int(cur != pre))

        dp = [int(1e20)] * (n + 1)
        dp[0] = 0
        for i in range(1, n + 1):
            for j in range(max(0, i - maxBoxes), i):
                weightSum = preWeight[i] - preWeight[j]
                if weightSum > maxWeight:
                    continue
                dp[i] = min(dp[i], dp[j] + preCost[i - 1] - preCost[j] + 2)
        return dp[-1]

    def boxDelivering2(self, boxes: List[List[int]], _: int, maxBoxes: int, maxWeight: int) -> int:
        n = len(boxes)
        preWeight = [0] + list(accumulate(box[1] for box in boxes))  # 重量前缀和
        preCost = [0]  # 运送次数前缀和 一次搬运前1个到前3个需要`preCost[2] - preCost[0]`次转移
        for (pre, _), (cur, _) in zip(boxes, boxes[1:]):
            preCost.append(preCost[-1] + int(cur != pre))

        dp = [int(1e20)] * (n + 1)
        dp[0] = 0
        queue = MonoQueue()
        for i in range(n + 1):
            while queue.minQueue and (
                i - queue.minQueue[0][1] > maxBoxes  # 超出数量限制
                or preWeight[i] - preWeight[queue.minQueue[0][1]] > maxWeight  # 超出重量限制
            ):
                queue.popleft()
            if queue:
                dp[i] = min(dp[i], preCost[i - 1] + queue.min + 2)
            if i < len(preCost):
                queue.append(dp[i] - preCost[i], i)
        return dp[-1]


print(Solution().boxDelivering2(boxes=[[1, 1], [2, 1], [1, 1]], _=2, maxBoxes=3, maxWeight=3))
# 输出：4
# 解释：最优策略如下：
# - 卡车将所有箱子装上车，到达码头 1 ，然后去码头 2 ，然后再回到码头 1 ，最后回到仓库，总共需要 4 趟行程。
# 所以总行程数为 4 。
# 注意到第一个和第三个箱子不能同时被卸货，因为箱子需要按顺序处理（也就是第二个箱子需要先被送到码头 2 ，然后才能处理第三个箱子）。

