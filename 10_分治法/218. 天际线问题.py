# 218. 天际线问题
# https://leetcode.cn/problems/the-skyline-problem/description/

from typing import List


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def getSkyline(self, buildings: List[List[int]]) -> List[List[int]]:
        """
        Divide and conquer approach:
        - Recursively compute skyline for left half and right half.
        - Merge two skylines in O(m+n) time.
        Overall O(n log n) time, O(n) space.
        """
        if not buildings:
            return []

        def merge(sky1: List[List[int]], sky2: List[List[int]]) -> List[List[int]]:
            i = j = 0
            h1 = h2 = 0
            merged: List[List[int]] = []

            # Merge by x-coordinate
            while i < len(sky1) and j < len(sky2):
                x1, y1 = sky1[i]
                x2, y2 = sky2[j]
                if x1 < x2:
                    x, h1 = x1, y1
                    i += 1
                elif x2 < x1:
                    x, h2 = x2, y2
                    j += 1
                else:  # x1 == x2
                    x = x1
                    h1 = y1
                    h2 = y2
                    i += 1
                    j += 1
                maxh = max2(h1, h2)
                # only add if height changed
                if not merged or merged[-1][1] != maxh:
                    merged.append([x, maxh])

            while i < len(sky1):
                x, y = sky1[i]
                if merged[-1][1] != y:
                    merged.append([x, y])
                i += 1
            while j < len(sky2):
                x, y = sky2[j]
                if merged[-1][1] != y:
                    merged.append([x, y])
                j += 1

            return merged

        def solve(l: int, r: int) -> List[List[int]]:
            if l == r:
                L, R, H = buildings[l]
                return [[L, H], [R, 0]]
            mid = (l + r) // 2
            left_sky = solve(l, mid)
            right_sky = solve(mid + 1, r)
            return merge(left_sky, right_sky)

        return solve(0, len(buildings) - 1)


if __name__ == "__main__":
    sol = Solution()
    buildings = [[2, 9, 10], [3, 7, 15], [5, 12, 12], [15, 20, 10], [19, 24, 8]]
    print(sol.getSkyline(buildings))
    # Expected: [[2,10],[3,15],[7,12],[12,0],[15,10],[20,8],[24,0]]
