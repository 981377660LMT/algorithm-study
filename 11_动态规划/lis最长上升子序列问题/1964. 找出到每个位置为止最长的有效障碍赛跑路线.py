from typing import List
from bisect import bisect_left, bisect_right

# 1 <= n <= 105
# 请你找出 obstacles 能构成的最长障碍路线的长度：
# 路线中每个障碍的高度都必须和前一个障碍 相同 或者 更高 。

# 即求每个位置结尾的最长上升子序列的长度
class Solution:
    def longestObstacleCourseAtEachPosition(self, obstacles: List[int]) -> List[int]:
        n = len(obstacles)
        res = [1] * n
        LIS = [obstacles[0]]

        # 可以取相等：使用bisectRight
        for i, num in list(enumerate(obstacles))[1:]:
            if num >= LIS[-1]:
                LIS.append(num)
                res[i] = len(LIS)
            else:
                index = bisect_right(LIS, num)
                res[i] = index + 1
                LIS[index] = num
        return res


print(Solution().longestObstacleCourseAtEachPosition(obstacles=[1, 2, 3, 2]))
# 输出：[1,2,3,3]
# 解释：每个位置的最长有效障碍路线是：
# - i = 0: [1], [1] 长度为 1
# - i = 1: [1,2], [1,2] 长度为 2
# - i = 2: [1,2,3], [1,2,3] 长度为 3
# - i = 3: [1,2,3,2], [1,2,2] 长度为 3
print(Solution().longestObstacleCourseAtEachPosition(obstacles=[2, 2, 1]))
# [1,2,1]
print(Solution().longestObstacleCourseAtEachPosition(obstacles=[5, 1, 5, 5, 1, 3, 4, 5, 1, 4]))
# [1,1,2,3,2,3,4,5,3,5]
