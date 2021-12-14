from typing import List

# 这个箱子被 顺时针旋转 90 度 ，由于重力原因，部分石头的位置会发生改变
# 每个石头会垂直掉落，直到它遇到障碍物，
# '#' 表示石头
# '*' 表示固定的障碍物
# '.' 表示空位置

# Shift stones then rotate
class Solution:
    def rotateTheBox(self, box: List[List[str]]) -> List[List[str]]:
        # 思路：遍历每行，把石头搬到movePosition
        for row in box:
            movePos = len(row) - 1
            for col in range(len(row) - 1, -1, -1):
                if row[col] == '*':
                    movePos = col - 1
                elif row[col] == '#':
                    row[col], row[movePos] = row[movePos], row[col]
                    movePos -= 1

        return list(zip(*box[::-1]))


print(Solution().rotateTheBox(box=[["#", ".", "*", "."], ["#", "#", "*", "."]]))
# 输出：[["#","."],
#       ["#","#"],
#       ["*","*"],
#       [".","."]]
