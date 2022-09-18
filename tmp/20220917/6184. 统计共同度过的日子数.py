# 所有日期的格式均为 "MM-DD" 。
# Alice 和 Bob 的到达日期都 早于或等于 他们的离开日期。
# 题目测试用例所给出的日期均为 非闰年 的有效日期。
# 请你返回 Alice和 Bob 同时在罗马的天数。
# [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31] 。


class Solution:
    def countDaysTogether(
        self, arriveAlice: str, leaveAlice: str, arriveBob: str, leaveBob: str
    ) -> int:
        """求区间相交长度/线段相交长度"""

        def mmdd2days(mmdd: str) -> int:
            mm, dd = map(int, mmdd.split("-"))
            return sum([31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31][: mm - 1]) + dd

        start1, end1, start2, end2 = map(mmdd2days, [arriveAlice, leaveAlice, arriveBob, leaveBob])
        return max(0, min(end1, end2) - max(start1, start2) + 1)
