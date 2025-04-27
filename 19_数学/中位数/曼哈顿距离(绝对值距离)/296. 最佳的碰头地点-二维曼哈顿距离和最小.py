# https://leetcode.cn/problems/best-meeting-point/description/
# !注意技巧：按照行和列分别扫描，天然有序

from typing import List


class Solution:
    def minTotalDistance(self, grid: List[List[int]]) -> int:
        """
        利用曼哈顿距离可分解到行和列两个一维问题的特性：
        在一维上，将所有点移动到中位数位置能最小化绝对距离之和。
        所以：
          1. 收集所有朋友的行坐标 rows（按升序），
             以及所有朋友的列坐标 cols（按升序）。
          2. 取 rows 和 cols 的中位数 hr, hc 作为最佳碰头点。
          3. 计算总距离：sum(|r-hr|) + sum(|c-hc|)。
        时间 O(m·n)，空间 O(k)，k 为朋友数。
        """
        m = len(grid)
        n = len(grid[0]) if m > 0 else 0

        # 收集所有行坐标（按行扫描，天然有序）
        rows: List[int] = []
        for i in range(m):
            for j in range(n):
                if grid[i][j] == 1:
                    rows.append(i)

        # 收集所有列坐标（按列扫描，天然有序）
        cols: List[int] = []
        for j in range(n):
            for i in range(m):
                if grid[i][j] == 1:
                    cols.append(j)

        mid = len(rows) // 2
        hr = rows[mid]
        hc = cols[mid]

        res = 0
        res += sum(abs(r - hr) for r in rows)
        res += sum(abs(c - hc) for c in cols)
        return res


if __name__ == "__main__":
    sol = Solution()
    grid1 = [[1, 0, 0, 0, 1], [0, 0, 0, 0, 0], [0, 0, 1, 0, 0]]
    # 三个朋友在 (0,0), (0,4), (2,2)
    # 最优碰头点在 (0,2) 或 (2,2) 等，距离和 = 6
    print(sol.minTotalDistance(grid1))  # 输出 6
