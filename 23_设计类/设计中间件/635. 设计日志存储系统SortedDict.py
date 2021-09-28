# TreeMap 提供了一个根据 key 的上界和下界生成一个子 Map 的功能，
# 刚好能够适用于题目的要求

from sortedcontainers import SortedDict
import datetime


class LogSystem:
    def __init__(self):
        self.info = SortedDict()

    # 给定日志的 id 和 timestamp ，将这个日志存入你的存储系统中
    def put(self, id: int, timestamp: str) -> None:
        time_arr = datetime.datetime.strptime(timestamp, "%Y:%m:%d:%H:%M:%S")
        if time_arr in self.info:
            self.info[time_arr].append(id)
        else:
            self.info[time_arr] = [id]

    # 返回在给定时间区间 [start, end] （包含两端）内的所有日志的 id
    def retrieve(self, s: str, e: str, granularity: str):
        start = datetime.datetime.strptime(s, "%Y:%m:%d:%H:%M:%S")
        end = datetime.datetime.strptime(e, "%Y:%m:%d:%H:%M:%S")

        if granularity == 'Year':
            start = datetime.datetime(year=start.year, month=1, day=1, hour=0, minute=0, second=0)
            end = datetime.datetime(year=end.year + 1, month=1, day=1, hour=0, minute=0, second=0)
        elif granularity == 'Month':
            start = datetime.datetime(
                year=start.year, month=start.month, day=1, hour=0, minute=0, second=0
            )
            if end.month == 12:
                end = datetime.datetime(
                    year=end.year + 1, month=1, day=1, hour=0, minute=0, second=0
                )
            else:
                end = datetime.datetime(
                    year=end.year, month=end.month + 1, day=1, hour=0, minute=0, second=0
                )
        elif granularity == 'Day':
            start = datetime.datetime(
                year=start.year, month=start.month, day=start.day, hour=0, minute=0, second=0
            )
            end = end + datetime.timedelta(days=1)
            end = datetime.datetime(
                year=end.year, month=end.month, day=end.day, hour=0, minute=0, second=0
            )
        elif granularity == 'Hour':
            start = datetime.datetime(
                year=start.year,
                month=start.month,
                day=start.day,
                hour=start.hour,
                minute=0,
                second=0,
            )
            end = end + datetime.timedelta(hours=1)
            end = datetime.datetime(
                year=end.year, month=end.month, day=end.day, hour=end.hour, minute=0, second=0
            )
        elif granularity == 'Minute':
            start = datetime.datetime(
                year=start.year,
                month=start.month,
                day=start.day,
                hour=start.hour,
                minute=start.minute,
                second=0,
            )
            end = end + datetime.timedelta(minutes=1)
            end = datetime.datetime(
                year=end.year,
                month=end.month,
                day=end.day,
                hour=end.hour,
                minute=end.minute,
                second=0,
            )
        elif granularity == 'Second':
            end = end + datetime.timedelta(seconds=1)

        res = []

        # 提取范围内的键
        keys = self.info.irange(start, end)
        for key in keys:
            if key >= end:
                break
            list = self.info[key]
            for id in list:
                res.append(id)

        return res


# Your LogSystem object will be instantiated and called as such:
# obj = LogSystem()
# obj.put(id,timestamp)
# param_2 = obj.retrieve(start,end,granularity)

