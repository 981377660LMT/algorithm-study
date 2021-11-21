from typing import List

# 学生们就近坐下 求总的最小移动距离

# ，一个萝卜一个坑
# 交换任意两个学生对应的座位不会产生更少的移动次数
class Solution:
    def minMovesToSeat(self, seats: List[int], students: List[int]) -> int:
        return sum(abs(a - b) for a, b in zip(sorted(seats), sorted(students)))


# 输入：seats = [3,1,5], students = [2,7,4]
# 输出：4
# 解释：学生移动方式如下：
# - 第一位学生从位置 2 移动到位置 1 ，移动 1 次。
# - 第二位学生从位置 7 移动到位置 5 ，移动 2 次。
# - 第三位学生从位置 4 移动到位置 3 ，移动 1 次。
# 总共 1 + 2 + 1 = 4 次移动。

