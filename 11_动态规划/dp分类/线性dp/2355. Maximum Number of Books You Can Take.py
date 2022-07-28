from typing import List

# https://leetcode.cn/problems/maximum-number-of-books-you-can-take/solution/python-dan-diao-zhan-d-by-981377660lmt-m03d/
# 母题:
# 1. 求等差数列的和
# 2. 单调栈

# 心路历程:
# !一开始用(index,是否选择前一项) 作为状态定义，后面发现这样不行
# !如何最暴力地计算每个位置能拿走的数量
# !直接往前遍历，等差数列求和
# 有的位置不满足怎么办
# !找出满足的一截等差数列 不满足的那一截之前算过


class Solution:
    def maximumBooks(self, books: List[int]) -> int:
        """
        n,nums[i]<=1e5
        从连续的书架上拿走最多的书 每次拿的书数量必须严格递增

        dp[i] 表示前i本书且取第i本书时 可以拿走的最大数量
        需要找到之前 arr[j] < arr[i] - (i-j) 的第一个j
        即对每个i找到左边第一个j 满足 arr[j] - j < arr[i] - i
        `dp[i]=dp[j]+(j+1到i这一段等差数列的和)`
        """
        n = len(books)
        firstJ = [-1] * n  # !对每个i找到左边第一个j 满足 arr[j] - j < arr[i] - i
        nums = [num - i for i, num in enumerate(books)]  # !对每个位置找到左边第一个比他严格小的数的位置 从右往左维护一个递增的单调栈
        stack = []
        for i in range(n - 1, -1, -1):
            while stack and nums[stack[-1]] > nums[i]:
                firstJ[stack.pop()] = i
            stack.append(i)

        dp = [0] * n
        dp[0] = books[0]
        for i in range(1, n):
            j = firstJ[i]
            count = min(i - j, books[i])  # 等差数列项数
            first, last = max(1, books[i] - count + 1), books[i]
            sum_ = (first + last) * count // 2  # j+1 到 i 均匀增长的一段的和
            dp[i] = sum_ + (dp[j] if j != -1 else 0)
        return max(dp)


# print(Solution().maximumBooks(books=[8, 2, 3, 7, 3, 4, 0, 1, 4, 3]))
# print(Solution().maximumBooks(books=[7, 0, 3, 4, 5]))
# print(Solution().maximumBooks(books=[8, 5, 2, 7, 9]))
