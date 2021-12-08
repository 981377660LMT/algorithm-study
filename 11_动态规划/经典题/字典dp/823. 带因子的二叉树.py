from typing import List

# 1 <= arr.length <= 1000
# arr 中的所有值 互不相同

# 用这些整数来构建二叉树，每个整数可以使用任意次数。
# 其中：每个非叶结点的值应等于它的两个子结点的值的乘积。
# 满足条件的二叉树一共有多少个？

# bottom up DP, smaller numbers are always the leaf nodes.
# 总结：排序+枚举根


class Solution:
    def numFactoredBinaryTrees(self, arr: List[int]) -> int:
        arr.sort()
        # 以root为根节点的二叉树个数
        dp = {root: 1 for root in arr}

        # 枚举根
        for index, root in enumerate(arr):
            for leafOne in arr[:index]:
                if root % leafOne == 0:
                    leafTwo = root // leafOne
                    if leafTwo in dp:
                        dp[root] += dp[leafOne] * dp[leafTwo]

        return sum(dp.values()) % int(1e9 + 7)


print(Solution().numFactoredBinaryTrees(arr=[2, 4, 5, 10]))
# 输出: 7
# 解释: 可以得到这些二叉树: [2], [4], [5], [10], [4, 2, 2], [10, 2, 5], [10, 5, 2].
