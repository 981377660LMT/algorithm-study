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
            fallTo = len(row) - 1
            for col in range(len(row) - 1, -1, -1):
                if row[col] == '*':
                    fallTo = col - 1
                elif row[col] == '#':
                    # 交换位置
                    row[col], row[fallTo] = row[fallTo], row[col]
                    fallTo -= 1

        return [list(col[::-1]) for col in zip(*box)]


print(Solution().rotateTheBox(box=[["#", ".", "*", "."], ["#", "#", "*", "."]]))
# 输出：[["#","."],
#       ["#","#"],
#       ["*","*"],
#       [".","."]]
