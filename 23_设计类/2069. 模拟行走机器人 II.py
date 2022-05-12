from typing import List

DIR = {
    0: "East",
    1: "North",
    2: "West",
    3: "South",
}


# 关键是初始化每个点的转向
class Robot:
    def __init__(self, width: int, height: int):
        self.isMoved = False
        self.idx = 0
        self.pos = []
        self.dirs = []

        for x in range(width):
            self.pos.append((x, 0))
            self.dirs.append(0)
        for y in range(1, height):
            self.pos.append((width - 1, y))
            self.dirs.append(1)
        for x in range(width - 2, -1, -1):
            self.pos.append((x, height - 1))
            self.dirs.append(2)
        for y in range(height - 2, 0, -1):
            self.pos.append((0, y))
            self.dirs.append(3)

        self.dirs[0] = 3

    def step(self, num: int) -> None:
        self.isMoved = True

        self.idx = (self.idx + num) % len(self.pos)

    def getPos(self) -> List[int]:

        return list(self.pos[self.idx])

    def getDir(self) -> str:
        if not self.isMoved:
            return "East"

        return DIR[self.dirs[self.idx]]


# Your Robot object will be instantiated and called as such:
# obj = Robot(width, height)
# obj.step(num)
# param_2 = obj.getPos()
# param_3 = obj.getDir()
