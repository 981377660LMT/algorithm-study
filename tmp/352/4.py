from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 现有 n 个机器人，编号从 1 开始，每个机器人包含在路线上的位置、健康度和移动方向。

# 给你下标从 0 开始的两个整数数组 positions、healths 和一个字符串 directions（directions[i] 为 'L' 表示 向左 或 'R' 表示 向右）。 positions 中的所有整数 互不相同 。

# 所有机器人以 相同速度 同时 沿给定方向在路线上移动。如果两个机器人移动到相同位置，则会发生 碰撞 。

# 如果两个机器人发生碰撞，则将 健康度较低 的机器人从路线中 移除 ，并且另一个机器人的健康度 减少 1 。幸存下来的机器人将会继续沿着与之前 相同 的方向前进。如果两个机器人的健康度相同，则将二者都从路线中移除。

# 请你确定全部碰撞后幸存下的所有机器人的 健康度 ，并按照原来机器人编号的顺序排列。即机器人 1 （如果幸存）的最终健康度，机器人 2 （如果幸存）的最终健康度等。 如果不存在幸存的机器人，则返回空数组。

# 在不再发生任何碰撞后，请你以数组形式，返回所有剩余机器人的健康度（按机器人输入中的编号顺序）。


# 注意：位置  positions 可能是乱序的。


# 定位到LR
# 模拟机器人碰撞，拆分函数
class Solution:
    def survivedRobotsHealths(
        self, positions: List[int], H: List[int], directions: str
    ) -> List[int]:
        # 函数不纯，不太好
        def collision(i: int, j: int) -> None:
            """编号为i和j的机器人碰撞-><-,改变res[i]和res[j]"""
            resI, resJ = remainHp[i], remainHp[j]
            if resI > resJ:
                remainHp[i] -= 1
                remainHp[j] = 0
            elif resI < resJ:
                remainHp[j] -= 1
                remainHp[i] = 0
            else:
                remainHp[i] = 0
                remainHp[j] = 0

        pairs: List[Tuple[int, str, int]] = [
            (rid, dir, pos) for rid, (dir, pos) in enumerate(zip(directions, positions))
        ]
        pairs.sort(key=lambda x: x[2])
        remainHp = H[:]
        leftToRight: List[int] = []
        for rightId, dir, _ in pairs:
            if dir == "R":
                leftToRight.append(rightId)
            else:  # <-
                if not leftToRight:
                    continue
                while leftToRight and remainHp[rightId]:
                    leftId = leftToRight[-1]
                    collision(leftId, rightId)
                    if not remainHp[leftId]:
                        leftToRight.pop()

        return [v for v in remainHp if v]


# positions = [3,5,2,6], healths = [10,10,15,12], directions = "RLRL"
print(Solution().survivedRobotsHealths([3, 5, 2, 6], [10, 10, 15, 12], "RLRL"))
