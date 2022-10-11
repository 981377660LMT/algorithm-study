"""给定二叉搜索树的插入顺序求深度"""
#  Depth of BST Given Insertion Order

# !每插入一个数，相当于划分区间
# https://leetcode.cn/problems/depth-of-bst-given-insertion-order/solution/tu-jie-si-lu-jian-dan-dai-ma-mei-cha-ru-30c8q/
# !新节点的最终父节点，一定是在原树中，与其 `绝对值之差最接近` 两个元素之一。
# !二分查找当前元素最接近的两个元素，取深度的较大值

from typing import List
from sortedcontainers import SortedDict


INF = int(1e20)


class Solution:
    def maxDepthBST(self, order: List[int]) -> int:
        """order 是从 1 到 n 的整数的一个排列。"""
        # 元素起点 => 高度
        depthMap = SortedDict()

        # 两个哨兵
        depthMap[-INF] = 0
        depthMap[INF] = 0

        depthMap[order[0]] = 1
        res = 1
        for num in order[1:]:
            pos = depthMap.bisect_right(num)
            lower, higher = pos - 1, pos
            lowerEntry, higherEntry = depthMap.peekitem(lower), depthMap.peekitem(higher)
            depth = max(lowerEntry[1], higherEntry[1]) + 1
            depthMap[num] = depth
            res = max(res, depth)

        return res


print(Solution().maxDepthBST(order=[2, 1, 4, 3]))
# Output: 3
# Explanation: The binary search tree has a depth of 3 with path 2->3->4.

dic = SortedDict({1: 2, 3: 0, 2: 6, 7: 9, 4: 33, 9: 88})

# irange是根据值的范围,返回切片
# islice与切片一样,要指定索引
print(dic, *dic.irange(3, 8), next(dic.islice(3, 5)))
