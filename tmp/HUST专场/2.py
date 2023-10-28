from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


# 老师在黑板上写出一个字符串 s 表示方程，该方程仅包含变量 x 、其对应系数和 '+' ， '-' 操作。
# 假设方程中 x 为 answer ，将 x 以字符串 "x=#answer" 的形式返回。
# 题目保证，如果方程中只有一个解，则 answer 的值是一个整数。
# 如果方程没有解或存在的解不为整数，请返回 "No solution" 。
# 如果方程有无限解，则返回 "Infinite solutions" 。


class Solution:
    def mathProblem(self, s: str) -> str:
        def split_line(line: str) -> Tuple[int, int]:
            if line[0] == "x":
                line = "1" + line
            text = line.replace("+x", "+1x").replace("-x", "-1x").replace("-", "+-").split("+")
            text = [seg for seg in text if len(seg) > 0]

            x = sum([int(t[:-1]) for t in text if t[-1] == "x"])  # coefficient
            num = sum([int(t) for t in text if t[-1] != "x"])  # constant
            return x, num

        line_left, line_right = s.split("=")
        x_left, num_left = split_line(line_left)
        x_right, num_right = split_line(line_right)
        x = x_left - x_right
        num = num_right - num_left

        if x == 0:
            if num == 0:
                return "Infinite solutions"
            else:
                return "No solution"
        else:
            return "x=" + str(num // x)
