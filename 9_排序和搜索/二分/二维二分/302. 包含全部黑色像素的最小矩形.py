from typing import List

# 给你一个大小为 m x n 的二进制矩阵 image 表示一张黑白图片，0 代表白色像素，1 代表黑色像素。
# 黑色像素相互连接，也就是说，图片中只会有一片连在一块儿的黑色像素。像素点是水平或竖直方向连接的。
# 给你两个整数 x 和 y 表示某一个黑色像素的位置，请你找出包含全部黑色像素的最小矩形（与坐标轴对齐），并返回该矩形的面积。
# 你必须设计并实现一个时间复杂度低于 O(mn) 的算法来解决此问题。


# 注意最左能力二分取l 最右能力二分取r
class Solution:
    def minArea(self, image: List[List[str]], x: int, y: int) -> int:
        if len(image) == 0 or len(image[0]) == 0:
            return 0

        m, n = len(image), len(image[0])

        up = self._cal_up(image, x, m, n)
        down = self._cal_down(image, x, m, n)
        left = self._cal_left(image, y, m, n)
        right = self._cal_right(image, y, m, n)

        return (down - up + 1) * (right - left + 1)

    def _cal_up(self, image, bound, m, n):
        l, r = 0, bound
        while l <= r:
            mid = (l + r) >> 1
            if "1" in image[mid]:
                r = mid - 1
            else:
                l = mid + 1
        up = l
        return up

    def _cal_down(self, image, bound, m, n):
        l, r = bound, m - 1
        while l <= r:
            mid = (l + r) >> 1
            if "1" in image[mid]:
                l = mid + 1
            else:
                r = mid - 1
        down = r
        return down

    def _cal_left(self, image, bound, m, n):
        l, r = 0, bound
        while l <= r:
            mid = (l + r) >> 1
            has_one = False
            for row in range(m):
                if image[row][mid] == "1":
                    has_one = True
                    break

            if has_one:
                r = mid - 1
            else:
                l = mid + 1
        left = l
        return left

    def _cal_right(self, image, bound, m, n):
        l, r = bound, n - 1
        while l <= r:
            mid = (l + r) >> 1
            has_one = False
            for row in range(m):
                if image[row][mid] == "1":
                    has_one = True
                    break
            if has_one:
                l = mid + 1
            else:
                r = mid - 1
        right = r
        return right


print(Solution().minArea([["0", "0", "1", "0"], ["0", "1", "1", "0"], ["0", "1", "0", "0"]], 0, 2))
