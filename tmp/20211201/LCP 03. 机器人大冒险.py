from typing import List, Set, Tuple


Point = Tuple[int, int]
# 小伙伴事先给机器人输入一串指令command，机器人就会无限循环这条指令的步骤进行移动。指令有两种：
# U: 向y轴正方向移动一格
# R: 向x轴正方向移动一格。
# 给定终点坐标(x, y)，返回机器人能否完好地到达终点


class Solution:
    def robot(self, command: str, obstacles: List[List[int]], x: int, y: int) -> bool:
        # 到达target需要的偏移量
        def calDiff(target: Point, cycleEnd: Point) -> Point:
            x_times = target[0] // cycleEnd[0]
            y_times = target[1] // cycleEnd[1]
            min_times = min(x_times, y_times)
            return (target[0] - cycleEnd[0] * min_times, target[1] - cycleEnd[1] * min_times)

        pos = (0, 0)
        path: Set[Point] = set([(0, 0)])
        for com in command:
            if com == 'U':
                pos = (pos[0], pos[1] + 1)
            else:
                pos = (pos[0] + 1, pos[1])
            path.add(pos)

        diff = calDiff((x, y), pos)

        if diff not in path:
            return False

        for ox, oy in obstacles:
            if ox <= x and oy <= y:
                if calDiff((ox, oy), pos) in path:
                    return False

        return True


print(Solution().robot(command="URR", obstacles=[], x=3, y=2))
# 输出：true
# 解释：U(0, 1) -> R(1, 1) -> R(2, 1) -> U(2, 2) -> R(3, 2)。
