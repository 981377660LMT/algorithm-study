from typing import List

# 1 <= rows <= 150
# 1 <= columns <= 150

# 请你返回有多少个 子矩形 的元素全部都是 1 。

# O(n ^ 2 * m)
# 用 left[i][j] 表示矩阵中位置为 (i, j) 的元素前面有多少个连续的 1，
# 然后遍历矩阵中的每一个元素，计算以这个元素`为右下角`的全 1 子矩阵有多少个。
# 这里以每一个遍历到的元素`为右下角`计算全 1 子矩阵，既不会遗漏也不会重复。


class Solution:
    def numSubmat(self, mat: List[List[int]]) -> int:
        m, n = len(mat), len(mat[0])

        countLeft = [[0] * n for _ in range(m)]
        for i in range(m):
            for j in range(n):
                if j == 0:
                    countLeft[i][j] = mat[i][j]
                else:
                    countLeft[i][j] = 0 if mat[i][j] == 0 else countLeft[i][j - 1] + 1

        # countUp = [[0] * n for _ in range(m)]
        # for i in range(m):
        #     for j in range(n):
        #         if i == 0:
        #             countUp[i][j] = mat[i][j]
        #         else:
        #             countUp[i][j] = 0 if mat[i][j] == 0 else countLeft[i - 1][j] + 1

        res = 0
        for i in range(m):
            for j in range(n):
                if mat[i][j] == 0:
                    continue
                minLen = countLeft[i][j]
                # 向上加大高度
                for upper in range(i, -1, -1):
                    minLen = min(minLen, countLeft[upper][j])
                    if minLen == 0:
                        break
                    res += minLen

        return res


print(Solution().numSubmat(mat=[[1, 0, 1], [1, 1, 0], [1, 1, 0]]))
# 输出：13
# 解释：
# 有 6 个 1x1 的矩形。
# 有 2 个 1x2 的矩形。
# 有 3 个 2x1 的矩形。
# 有 1 个 2x2 的矩形。
# 有 1 个 3x1 的矩形。
# 矩形数目总共 = 6 + 2 + 3 + 1 + 1 = 13 。

# mat = [[1,0,1],
#        [1,1,0],
#        [1,1,0]]
# becomes
# mat = [[1,0,1],
#        [2,1,0],
#        [3,2,0]]
