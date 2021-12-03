from typing import List

# 你需要用这一辆车把一些箱子从仓库运送到码头。这辆卡车每次运输有 箱子数目的限制 和 总重量的限制 。
# 箱子需要按照 数组顺序 运输
# 对于在卡车上的箱子，我们需要 按顺序 处理它们

# 请你返回将所有箱子送到相应码头的 最少行程 次数。

# https://leetcode-cn.com/problems/delivering-boxes-from-storage-to-ports/solution/zai-liang-chong-tan-xin-ce-lue-zhong-zuo-qepx/
class Solution:
    def boxDelivering(
        self, boxes: List[List[int]], portsCount: int, maxBoxes: int, maxWeight: int
    ) -> int:
        ...


print(
    Solution().boxDelivering(boxes=[[1, 1], [2, 1], [1, 1]], portsCount=2, maxBoxes=3, maxWeight=3)
)
# 输出：4
# 解释：最优策略如下：
# - 卡车将所有箱子装上车，到达码头 1 ，然后去码头 2 ，然后再回到码头 1 ，最后回到仓库，总共需要 4 趟行程。
# 所以总行程数为 4 。
# 注意到第一个和第三个箱子不能同时被卸货，因为箱子需要按顺序处理（也就是第二个箱子需要先被送到码头 2 ，然后才能处理第三个箱子）。


# 太难了 放弃
