from typing import List


class Solution:
    def arrayChange(self, nums: List[int], operations: List[List[int]]) -> List[int]:
        """ !预处理opt的思想  用pre和next字典记录前驱后继 可以处理重复元素
        
        注意和并查集的区别
        并查集会不断找根,而next/pre只能连接着传递相邻关系
        """

        # !反向倒着修改的技巧 这种做法的好处在于，
        # !对于 operations 靠前的那些 x 到 y 的映射，我们是知道 x 最终要映射到哪个数字的。
        # !最终所有的点的属性只跟对其最后一次更新有关，所以可以选择逆序模拟
        next = dict()
        for a, b in operations[::-1]:
            next[a] = next.get(b, b)
        return [next.get(v, v) for v in nums]


print(Solution().arrayChange(nums=[1, 2, 4, 6, 1], operations=[[1, 3], [4, 7], [6, 1]]))
print(
    Solution().arrayChange(
        nums=[1, 2, 4, 6, 1], operations=[[1, 100], [4, 7], [6, 1], [1, 0], [100, 2]]
    )
)

