from typing import List


class Robot:

    TO_DIR = {
        0: "East",
        1: "North",
        2: "West",
        3: "South",
    }

    def __init__(self, width: int, height: int):
        self.moved = False
        self.idx = 0
        self.pos = list()
        self.dirs = list()

        pos_, dirs_ = self.pos, self.dirs

        for i in range(width):
            pos_.append((i, 0))
            dirs_.append(0)
        for i in range(1, height):
            pos_.append((width - 1, i))
            dirs_.append(1)
        for i in range(width - 2, -1, -1):
            pos_.append((i, height - 1))
            dirs_.append(2)
        for i in range(height - 2, 0, -1):
            pos_.append((0, i))
            dirs_.append(3)

        dirs_[0] = 3

    def move(self, num: int) -> None:
        self.moved = True
        self.idx = (self.idx + num) % len(self.pos)

    def getPos(self) -> List[int]:
        return list(self.pos[self.idx])

    def getDir(self) -> str:
        if not self.moved:
            return "East"
        return Robot.TO_DIR[self.dirs[self.idx]]


# Your Robot object will be instantiated and called as such:
# obj = Robot(width, height)
# obj.step(num)
# param_2 = obj.getPos()
# param_3 = obj.getDir()
