from typing import List


class Solution:
    def analysisHistogram(self, heights: List[int], cnt: int) -> List[int]:
        heights = sorted(heights)
        res = []
        minDiff = int(1e18)
        for right in range(cnt - 1, len(heights)):
            curDiff = heights[right] - heights[right - (cnt - 1)]
            if curDiff < minDiff:
                minDiff = curDiff
                # 注意边界:
                # 要取到right，所以切片右边界为right+1
                # 切片长cnt，所以切片左边界为right+1-cnt
                res = heights[right + 1 - cnt : right + 1]
        return res


print(Solution().analysisHistogram(heights=[3, 2, 7, 6, 1, 8], cnt=3))
print(Solution().analysisHistogram(heights=[4, 6, 1, 8, 4, 10], cnt=4))
