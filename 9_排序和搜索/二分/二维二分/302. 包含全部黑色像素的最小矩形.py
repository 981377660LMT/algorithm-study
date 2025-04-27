from typing import List

# 给你一个大小为 m x n 的二进制矩阵 image 表示一张黑白图片，0 代表白色像素，1 代表黑色像素。
# 黑色像素相互连接，也就是说，图片中只会有一片连在一块儿的黑色像素。像素点是水平或竖直方向连接的。
# 给你两个整数 x 和 y 表示某一个黑色像素的位置，请你找出包含全部黑色像素的最小矩形（与坐标轴对齐），并返回该矩形的面积。
# 你必须设计并实现一个时间复杂度低于 O(mn) 的算法来解决此问题。
#
# 注意最左能力二分取l 最右能力二分取r

from typing import List


class Solution:
    def minArea(self, image: List[List[str]], x: int, y: int) -> int:
        """
        Binary search on rows and columns to find the minimal enclosing rectangle.
        Time: O(n log m + m log n), better than O(mn).
        """
        m, n = len(image), len(image[0])

        def has_black_row(r: int) -> bool:
            return "1" in image[r]

        def has_black_col(c: int) -> bool:
            for i in range(m):
                if image[i][c] == "1":
                    return True
            return False

        # Generic binary search: find first index in [lo, hi) where check(idx) is True
        def bs(lo: int, hi: int, check) -> int:
            while lo < hi:
                mid = (lo + hi) // 2
                if check(mid):
                    hi = mid
                else:
                    lo = mid + 1
            return lo

        top = bs(0, x, has_black_row)
        bottom = bs(x + 1, m, lambda r: not has_black_row(r)) - 1

        left = bs(0, y, has_black_col)
        right = bs(y + 1, n, lambda c: not has_black_col(c)) - 1

        return (bottom - top + 1) * (right - left + 1)


if __name__ == "__main__":
    image = ["0010", "0110", "0100"]

    grid = [list(row) for row in image]
    sol = Solution()
    # 已知一个黑色像素在 (0,2)
    print(sol.minArea(grid, 0, 2))  # 输出 6
