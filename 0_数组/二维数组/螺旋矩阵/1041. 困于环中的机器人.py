# 正方向参考数轴xy
class Solution:
    def isRobotBounded(self, instructions: str) -> bool:
        x, y, dx, dy = 0, 0, 0, 1
        for i in instructions:
            if i == 'R':
                dx, dy = dy, -dx
            if i == 'L':
                dx, dy = -dy, dx
            if i == 'G':
                x, y = x + dx, y + dy
        return (x, y) == (0, 0) or (dx, dy) != (0, 1)

    def isRobotBounded2(self, ins: str) -> bool:
        direction = complex(0, 1)
        position = complex(0, 0)
        for c in ins:
            if c == 'L':
                direction *= complex(0, 1)
            elif c == 'R':
                direction *= complex(0, -1)
            else:
                position += direction
        return position == complex(0, 0) or direction != complex(0, 1)


# 只有在平面中存在环使得机器人永远无法离开时，返回 true。否则，返回 false。
# 1. 走到原点
# 2. 不在原点但是方向变化了：四次迭代就能变回初始方向
