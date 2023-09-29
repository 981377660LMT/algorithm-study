# 2852. 所有单元格的远离程度之和
# https://leetcode.cn/problems/sum-of-remoteness-of-all-cells/description/
# 序号0起始的二维数组grid，每个正整数数表示一个区域，-1表示阻挡。
# 你可以从一个区域自由移动到任何相邻的区域。
# 对于任意点(i,j)，其遥远度定义
# 1）若为区域，则为整个矩阵内所有无法从该点到达的区域之和；
# 2）若为阻挡则为0。
# 返回整个矩阵的遥远度总和。
# n<=300

# 计算每个联通分量的和


from collections import defaultdict
from typing import List

from UnionFind import UnionFindArray


class Solution:
    def sumRemoteness(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        uf = UnionFindArray(ROW * COL)
        for r in range(ROW):
            for c in range(COL):
                if grid[r][c] == -1:
                    continue
                if r > 0 and grid[r - 1][c] != -1:
                    uf.union(r * COL + c, (r - 1) * COL + c)
                if c > 0 and grid[r][c - 1] != -1:
                    uf.union(r * COL + c, r * COL + c - 1)

        groups = uf.getGroups()
        groupSum = defaultdict(int)
        for root, members in groups.items():
            if grid[root // COL][root % COL] == -1:
                continue
            for member in members:
                groupSum[root] += grid[member // COL][member % COL]
        allSum = sum(groupSum.values())

        res = 0
        for r in range(ROW):
            for c in range(COL):
                if grid[r][c] != -1:
                    res += allSum - groupSum[uf.find(r * COL + c)]
        return res


if __name__ == "__main__":
    print(Solution().sumRemoteness([[-1, 1, -1], [5, -1, 4], [-1, 3, -1]]))
