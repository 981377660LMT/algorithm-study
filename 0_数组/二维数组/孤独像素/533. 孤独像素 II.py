from typing import List

# 行 r 和列 c 中的黑色像素恰好有 target 个。
# 符合规则一的每一行的元素都要完全相同。
class Solution:
    def findBlackPixel(self, picture: List[List[str]], target: int) -> int:
        res, m, n = 0, len(picture), len(picture[0])
        rows = [''.join(row) for row in picture]
        for c in range(n):
            # 每一行的状态
            rowState = [rows[r] for r in range(m) if picture[r][c] == 'B']
            if len(rowState) == target == rowState[0].count('B') and len(set(rowState)) == 1:
                res += target
        return res


print(
    Solution().findBlackPixel(
        picture=[
            ["W", "B", "W", "B", "B", "W"],
            ["W", "B", "W", "B", "B", "W"],
            ["W", "B", "W", "B", "B", "W"],
            ["W", "W", "B", "W", "B", "W"],
        ],
        target=3,
    )
)

