from typing import Tuple

# 思路：补全系数，一切都是加
# 1. +-x全部变成+1x -1x
# 2. - 变成 +-
# 3. 按照加号分割所有
class Solution:
    def solveEquation(self, equation: str) -> str:
        def split_line(line: str) -> Tuple[int, int]:
            if line[0] == 'x':
                line = '1' + line
            text = line.replace('+x', '+1x').replace('-x', '-1x').replace('-', '+-').split('+')
            text = [i for i in text if len(i) > 0]

            x = sum([int(t[:-1]) for t in text if t[-1] == 'x'])  # coefficient
            num = sum([int(t) for t in text if t[-1] != 'x'])  # constant
            return x, num

        line_left, line_right = equation.split('=')
        x_left, num_left = split_line(line_left)
        x_right, num_right = split_line(line_right)
        x = x_left - x_right
        num = num_right - num_left

        if x == 0:
            if num == 0:
                return 'Infinite solutions'
            else:
                return 'No solution'
        else:
            return 'x=' + str(num // x)


print(Solution().solveEquation("x+5-3+x=6+x-2"))
