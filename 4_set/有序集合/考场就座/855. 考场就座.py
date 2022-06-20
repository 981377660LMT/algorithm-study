from bisect import insort_left

# 当学生进入考场后，他必须坐在能够使他与离他最近的人之间的距离达到最大化的座位上。
# https://leetcode-cn.com/problems/exam-room/solution/kao-chang-jiu-zuo-by-leetcode/

# !注意这道题遍历查询索引次数很多 因此SortedList并不划算


class ExamRoom:
    def __init__(self, n: int):
        self.cap = n
        self.seats = []

    # 坐在编号最小的座位上,离他最近的人之间的距离达到最大化的座位上
    def seat(self) -> int:
        if not self.seats:
            cand = 0
        else:
            # 遍历座位找最长间隔
            resMax, cand = self.seats[0], 0
            for i, cur in enumerate(self.seats):
                if i == 0:
                    continue
                pre = self.seats[i - 1]
                curMax = (cur - pre) >> 1  # !离他最近的人之间的距离达到最大化且编号最小
                if curMax > resMax:
                    resMax, cand = curMax, pre + curMax

            # 考虑坐在最后一个位置
            lastDist = self.cap - 1 - self.seats[-1]
            if lastDist > resMax:
                cand = self.cap - 1

        insort_left(self.seats, cand)
        return cand

    # 坐在座位 p 上的学生现在离开了
    # 每次调用 ExamRoom.leave(p) 时都保证有学生坐在座位 p 上
    def leave(self, p: int) -> None:
        self.seats.remove(p)


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

