#  Depth of BST Given Insertion Order
from typing import List
from sortedcontainers import SortedDict

INF = 0x3F3F3F3F

# 二分查找当前元素最接近的两个元素，取深度的较大值


class Solution:
    def maxDepthBST(self, order: List[int]) -> int:
        # 元素=>高度
        sd = SortedDict()
        # 初始化边界
        sd[0] = 0
        sd[INF] = 0
        sd[order[0]] = 1

        res = 1
        for num in order[1:]:
            upper = sd.bisect_right(num)
            depth = max(sd.peekitem(upper)[1], sd.peekitem(upper - 1)[1]) + 1
            sd[num] = depth
            res = max(res, depth)

        return res


print(Solution().maxDepthBST(order=[2, 1, 4, 3]))
# Output: 3
# Explanation: The binary search tree has a depth of 3 with path 2->3->4.

dic = SortedDict({1: 2, 3: 0, 2: 6, 7: 9, 4: 33, 9: 88})

# irange是根据值的范围,返回切片
# islice与切片一样,要指定索引
print(dic, *dic.irange(3, 8), next(dic.islice(3, 5)))

