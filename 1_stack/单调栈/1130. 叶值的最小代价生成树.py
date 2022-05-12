from typing import List

# 每个节点都有 0 个或是 2 个子节点。
# 数组 arr 中的值与树的中序遍历中每个`叶节点`的值一一对应
# 每个非叶节点的值等于其左子树和右子树中`叶节点的最大值的乘积`。
# 在所有这样的二叉树中，返回每个非叶节点的值的最小可能总和


# 这个题实际上是将数组中相邻的数两两合并，计算他们的乘积之和，
# 求最小的乘积之和。合并相邻的两个数之后得到的是较大的一个数。

# 想让 mct 值最小，那么值较小的叶子节点就要尽量放到底部，
# 值较大的叶子节点要尽量放到靠上的部分。
# `因为越是底部的叶子节点，被用来做乘法的次数越多`。

# 1.每次pop`最小数`，然后乘以相邻的值
# 为防止 IndexError: list index out of range 错误，我们使用切片相加
# 优化：这个找最小数的过程可以用单调栈代替

# 2.不断寻找极小值。通过维护一个单调递减栈就可以找到一个极小值
# 栈顶存的一直是当时能找到的最小值
# https://leetcode-cn.com/problems/minimum-cost-tree-from-leaf-values/solution/wei-shi-yao-dan-diao-di-jian-zhan-de-suan-fa-ke-xi/


class Solution:
    def mctFromLeafValues1(self, arr: List[int]) -> int:
        res = 0
        for _ in range(len(arr) - 1):
            i = arr.index(min(arr))
            res += min(arr[i - 1 : i] + arr[i + 1 : i + 2]) * arr.pop(i)
        return res

    def mctFromLeafValues2(self, arr: List[int]) -> int:
        """每个数作为最小值"""
        res = 0
        stack = [int(1e20)]
        for num in arr:
            while stack and stack[-1] < num:
                min_ = stack.pop()
                res += min(stack[-1], num) * min_
            stack.append(num)

        while len(stack) > 2:
            res += stack.pop() * stack[-1]
        return res


print(Solution().mctFromLeafValues2(arr=[6, 2, 4]))
# 输出：32
# 解释：
# 有两种可能的树，第一种的非叶节点的总和为 36，第二种非叶节点的总和为 32。

#     24            24
#    /  \          /  \
#   12   4        6    8
#  /  \               / \
# 6    2             2   4


# 2*6+6*4
# 2*4+4*6
