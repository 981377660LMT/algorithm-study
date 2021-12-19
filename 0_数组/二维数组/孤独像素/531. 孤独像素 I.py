from typing import List

# 请你统计并返回图像中 黑色 孤独像素的数量。
class Solution:
    def findLonelyPixel(self, picture: List[List[str]]) -> int:
        m, n = len(picture), len(picture[0])
        rows = [0] * m
        cols = [0] * n
        for i in range(m):
            for j in range(n):
                if picture[i][j] == 'B':
                    rows[i] += 1
                    cols[j] += 1

        res = 0
        for i in range(m):
            for j in range(n):
                if picture[i][j] == 'B' and rows[i] == 1 and cols[j] == 1:
                    res += 1
        return res


print(Solution().findLonelyPixel(picture=[["W", "W", "B"], ["W", "B", "W"], ["B", "W", "W"]]))
