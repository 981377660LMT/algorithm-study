from typing import List

# peek finding
# 1 <= mat[i][j] <= 105
# 任意两个相邻元素均不相等.
# 一个 2D 网格中的 顶峰元素 是指那些 严格大于 其相邻格子(上、下、左、右)的元素。
# 找出 任意一个 顶峰元素 mat[i][j] 并 返回其位置 [i,j]


# 相邻的三行中，中间的那一行的最大值，如果大于另外两行的最大值，则中间行的最大值的对应元素必定为顶峰
# 如果三行的最大值中，不是中间行的最大，那搜索范围就往最大的那一行的方向缩


class Solution:
    def findPeakGrid(self, mat: List[List[int]]) -> List[int]:
        m, _ = len(mat), len(mat[0])
        up, down = 0, m - 1
        while up <= down:
            mid = (up + down) >> 1
            midVal = max(mat[mid])
            upVal = max(mat[mid - 1]) if mid > 0 else -1
            downVal = max(mat[mid + 1]) if mid + 1 < m else -1

            if midVal > max(upVal, downVal):
                return [mid, mat[mid].index(midVal)]
            elif downVal > max(midVal, upVal):
                up = mid + 1
            else:
                down = mid - 1

        return [up, mat[up].index(max(mat[up]))]


print(Solution().findPeakGrid(mat=[[10, 20, 15], [21, 30, 14], [7, 16, 32]]))

# 输出: [1,1]
# 解释: 30和32都是顶峰元素，所以[1,1]和[2,2]都是可接受的答案。
