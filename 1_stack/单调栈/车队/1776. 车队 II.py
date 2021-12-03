from typing import List

# 其中 cars[i] = [positioni, speedi]
# 题目保证 positioni < positioni+1 。

# 一旦两辆车相遇，它们会合并成一个车队
# 速度为这个车队里 最慢 一辆车的速度。
# 返回一个数组 answer ，其中 answer[i] 是第 i 辆车与下一辆车相遇的时间（单位：秒），
# 如果这辆车不会与下一辆车相遇，则 answer[i] 为 -1


# 左侧的点，速度再大，都会被它右侧的点拦住。
# 一个点能不能追上右侧的点，与它左侧的点无关
# 从`右往左遍历即可`


class Solution:
    def getCollisionTimes(self, cars: List[List[int]]) -> List[float]:
        res = []
        # 右到左遍历,栈顶保存前面的车的信息(position, speed, collideTime)
        stack = []

        for position, speed in reversed(cars):
            # 车速慢，不会撞到前面一部车，弹出
            # 撞车之前前面车已经撞了，弹出
            while stack and (
                (speed <= stack[-1][1])
                or ((stack[-1][0] - position) / (speed - stack[-1][1]) >= stack[-1][2])
            ):
                stack.pop()
            if not stack:
                stack.append((position, speed, 0x7FFFFFFF))
                res.append(-1)
            else:
                # 这辆车与右边的车会相撞
                collideTime = (stack[-1][0] - position) / (speed - stack[-1][1])
                stack.append((position, speed, collideTime))
                res.append(collideTime)

        res.reverse()
        return res


print(Solution().getCollisionTimes(cars=[[1, 2], [2, 1], [4, 3], [7, 2]]))
# 输出：[1.00000,-1.00000,3.00000,-1.00000]
# 解释：经过恰好 1 秒以后，第一辆车会与第二辆车相遇，并形成一个 1 m/s 的车队。经过恰好 3 秒以后，第三辆车会与第四辆车相遇，并形成一个 2 m/s 的车队。
