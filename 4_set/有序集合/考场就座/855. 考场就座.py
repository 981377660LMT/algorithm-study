from sortedcontainers import SortedList

# 当学生进入考场后，他必须坐在能够使他与离他最近的人之间的距离达到最大化的座位上。
# https://leetcode-cn.com/problems/exam-room/solution/kao-chang-jiu-zuo-by-leetcode/


class ExamRoom:
    def __init__(self, n: int):
        self.capacity = n
        self.students = SortedList()

    # 坐在编号最小的座位上,离他最近的人之间的距离达到最大化的座位上
    def seat(self) -> int:
        if not self.students:
            choice = 0
        else:
            # 遍历座位找最长间隔
            maxDist, choice = self.students[0], 0
            for i, cur in enumerate(self.students):
                if i == 0:
                    continue
                pre = self.students[i - 1]
                curMaxDist = (cur - pre) >> 1
                if curMaxDist > maxDist:
                    maxDist, choice = curMaxDist, pre + curMaxDist

            # 考虑坐在最后一个位置
            lastDist = self.capacity - 1 - self.students[-1]
            if lastDist > maxDist:
                choice = self.capacity - 1

        self.students.add(choice)
        return choice

    # 坐在座位 p 上的学生现在离开了
    # 每次调用 ExamRoom.leave(p) 时都保证有学生坐在座位 p 上
    def leave(self, p: int) -> None:
        self.students.discard(p)


# ["ExamRoom","seat","seat","seat","seat","leave","seat"], [[10],[],[],[],[],[4],[]]
# 输出：[null,0,9,4,2,null,5]
# 解释：
# ExamRoom(10) -> null
# seat() -> 0，没有人在考场里，那么学生坐在 0 号座位上。
# seat() -> 9，学生最后坐在 9 号座位上。
# seat() -> 4，学生最后坐在 4 号座位上。
# seat() -> 2，学生最后坐在 2 号座位上。
# leave(4) -> null
# seat() -> 5，学生最后坐在 5 号座位上。

room = ExamRoom(10)
room.seat()
room.seat()
room.seat()
room.seat()
print(room.__dict__)

