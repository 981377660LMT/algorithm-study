from typing import List


# !假设没有任何限制 需要预处理operations


class Solution:
    def arrayChange(self, nums: List[int], operations: List[List[int]]) -> List[int]:
        """查找+修改 很自然的思想
        
        注意查找的indexOf 复杂度需要O(1) 不能是O(n)

        无法处理重复元素
        """
        # O(n)的查找 JS可以过
        # for a, b in operations:
        #     index = nums.index(a)
        #     nums[index] = b
        # return nums

        # O(1)的查找
        indexMap = {v: i for i, v in enumerate(nums)}
        for a, b in operations:
            index = indexMap[a]
            indexMap.pop(a)  # !删除是pop 也可以del 但是pop可读性更好
            nums[index] = b
            indexMap[b] = index
        return nums

    def arrayChange2(self, nums: List[int], operations: List[List[int]]) -> List[int]:
        """ !预处理opt的思想  用pre和next字典记录前驱后继 可以处理重复元素
        
        注意和并查集的区别
        并查集会不断找根,而next/pre只能连接着传递相邻关系
        """
        # pre, next = dict(), dict()
        # for a, b in operations:
        #     pre[b] = pre.get(a, a)  # !感受到了dict.get的好处
        #     next[pre[b]] = b
        # return [next.get(v, v) for v in nums]

        # !反向倒着修改的技巧 这种做法的好处在于，
        # !对于 operations 靠前的那些 x 到 y 的映射，我们是知道 x 最终要映射到哪个数字的。
        # !最终所有的点的属性只跟对其最后一次更新有关，所以可以选择逆序模拟
        next = dict()
        for a, b in operations[::-1]:
            next[a] = next.get(b, b)
        return [next.get(v, v) for v in nums]


print(Solution().arrayChange(nums=[1, 2, 4, 6, 1], operations=[[1, 3], [4, 7], [6, 1]]))
print(
    Solution().arrayChange2(
        nums=[1, 2, 4, 6, 1], operations=[[1, 100], [4, 7], [6, 1], [1, 0], [100, 2]]
    )
)

