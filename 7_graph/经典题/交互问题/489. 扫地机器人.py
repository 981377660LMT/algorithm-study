class Robot:
    def move(self):
        """
       Returns true if the cell in front is open and robot moves into the cell.
       Returns false if the cell in front is blocked and robot stays in the current cell.
       :rtype bool
       """

    def turnLeft(self):
        """
       Robot will stay in the same cell after calling turnLeft/turnRight.
       Each turn will be 90 degrees.
       :rtype void
       """

    def turnRight(self):
        """
       Robot will stay in the same cell after calling turnLeft/turnRight.
       Each turn will be 90 degrees.
       :rtype void
       """

    def clean(self):
        """
       Clean the current cell.
       :rtype void
       """


# 输入只用于初始化房间和机器人的位置。你需要“盲解”这个问题
# 时间复杂度：O(4^N−M)，其中 N 是房间的大小，M 是障碍物的数量。
class Solution:
    def cleanRoom(self, robot: Robot):
        """
        :type robot: Robot
        :rtype: None
        """

        def go_back() -> None:
            robot.turnRight()
            robot.turnRight()
            robot.move()
            robot.turnRight()
            robot.turnRight()

        def bt(row: int, col: int, di: int):
            visited.add((row, col))
            robot.clean()
            # going clockwise : 0: 'up', 1: 'right', 2: 'down', 3: 'left'
            for i in range(4):
                ndi = (i + di) % 4
                dr, dc = dirs[ndi]
                nr = row + dr
                nc = col + dc
                if (nr, nc) not in visited and robot.move():
                    bt(nr, nc, ndi)
                    # 回溯
                    go_back()
                robot.turnRight()

        # going clockwise : 0: 'up', 1: 'right', 2: 'down', 3: 'left'
        dirs = [(-1, 0), (0, 1), (1, 0), (0, -1)]
        visited = set()
        bt(0, 0, 0)


# 输入:
# room = [
#   [1,1,1,1,1,0,1,1],
#   [1,1,1,1,1,0,1,1],
#   [1,0,1,1,1,1,1,1],
#   [0,0,0,1,0,0,0,0],
#   [1,1,1,1,1,1,1,1]
# ],
# row = 1,
# col = 3

# 解析:
# 房间格栅用0或1填充。0表示障碍物，1表示可以通过。
# 机器人从row=1，col=3的初始位置出发。在左上角的一行以下，三列以右。

