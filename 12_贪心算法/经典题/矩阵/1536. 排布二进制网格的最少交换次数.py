from typing import List

# 给你一个 n x n 的二进制网格 grid，每一次操作中，你可以选择网格的 相邻两行 进行交换。
# 一个符合要求的网格需要满足主对角线以上的格子全部都是 0 。
# 主对角线指的是从 (1, 1) 到 (n, n) 的这些格子。
# 请你返回使网格满足要求的`最少操作次数`，如果无法使网格符合要求，请你返回 -1 。

# 思路:预处理每行右边有多少个连续的0;
# 最少操作次数:找到第一个满足条件的元素，模拟移动到第一行


class Solution:
    def minSwaps(self, grid: List[List[int]]) -> int:
        row, col = len(grid), len(grid[0])
        rightMost = [0] * row
        for r in range(row):
            # 最右边的1,没有1则为-1
            rightMost[r] = next((c for c in reversed(range(col)) if grid[r][c] == 1), -1)
        res = 0
        for r in range(row):
            for i, oneIndex in enumerate(rightMost):
                if oneIndex <= r:
                    res += i
                    rightMost.pop(i)
                    break
            else:
                # for里没break,执行else:没找到
                return -1
        return res


print(Solution().minSwaps(grid=[[0, 0, 1], [1, 1, 0], [1, 0, 0]]))
# 3
