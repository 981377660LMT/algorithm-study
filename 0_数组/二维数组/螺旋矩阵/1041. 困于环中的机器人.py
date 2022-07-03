# 只有在平面中存在环使得机器人永远无法离开时，返回 true。否则，返回 false。
# !1. 走到原点
# !2. 不在原点但是方向变化了：四次迭代就能变回初始方向

# up/right/down/left
DIR4 = ((-1, 0), (0, 1), (1, 0), (0, -1))


class Solution:
    def isRobotBounded(self, instructions: str) -> bool:
        index = 0
        r, c = 0, 0
        for i in instructions:
            if i == "R":
                index = (index + 1) % 4
            if i == "L":
                index = (index - 1) % 4
            if i == "G":
                dr, dc = DIR4[index]
                r, c = r + dr, c + dc
        return (r, c) == (0, 0) or index != 0

    def isRobotBounded2(self, ins: str) -> bool:
        """复数方向"""
        direction = complex(0, 1)
        position = complex(0, 0)
        for c in ins:
            if c == "L":
                direction *= complex(0, 1)
            elif c == "R":
                direction *= complex(0, -1)
            else:
                position += direction
        return position == complex(0, 0) or direction != complex(0, 1)
