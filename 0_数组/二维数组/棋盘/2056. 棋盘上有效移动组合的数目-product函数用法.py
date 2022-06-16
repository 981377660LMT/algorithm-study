# product穷举模拟
# 1 <= n <= 4

from itertools import product
from typing import Generator, List, Literal, Tuple


DIRS = {
    "rook": ((0, 1), (0, -1), (1, 0), (-1, 0)),
    "bishop": ((1, 1), (1, -1), (-1, 1), (-1, -1)),
    "queen": ((0, 1), (0, -1), (1, 0), (-1, 0), (1, 1), (1, -1), (-1, 1), (-1, -1)),
}

Role = Literal["rook", "bishop", "queen"]
State = Tuple[int, int, int, int, int]

# 棋盘下标从 1 开始
# python里的product函数表示笛卡尔积，非常适合**穷举所有对象的状态的组合**
# 一般product会配合数组解构使用 product(*(iterable for _ in range(n)))


class Solution:
    def countCombinations(self, pieces: List[Role], positions: List[List[int]]) -> int:
        def getNextState(role: Role, curX: int, curY: int) -> Generator[State, None, None]:
            """每一个状态用起点，方向，距离表示"""
            yield curX, curY, 0, 0, 0
            for dist, (dx, dy) in product(range(1, 8), DIRS[role]):
                if 1 <= curX + dx * dist <= 8 and 1 <= curY + dy * dist <= 8:
                    yield curX, curY, dx, dy, dist

        def checkAllStates(allStates: Tuple[State, ...]) -> bool:
            """所有棋子同时移动，检查是否有冲突"""
            for dist1 in range(8):
                visited = set()
                for x, y, dx, dy, dist2 in allStates:
                    curDist = min(dist1, dist2)
                    if (x + dx * curDist, y + dy * curDist) in visited:
                        return False
                    visited.add((x + dx * curDist, y + dy * curDist))
            return True

        return sum(
            checkAllStates(allStates)
            for allStates in product(
                *(getNextState(role, x, y) for role, (x, y) in zip(pieces, positions))
            )
        )


print(Solution().countCombinations(pieces=["queen", "bishop"], positions=[[5, 7], [3, 4]]))
