from typing import Tuple


class DateTime:
    _MONTH_DAYS = (0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31)

    @staticmethod
    def from_int(x: int) -> "DateTime":
        y = x * 400 // 146097 + 1
        d = x - DateTime(y, 1, 1).to_int()
        m = 1
        while d >= 28:
            k = DateTime._MONTH_DAYS[m] + (m == 2 and DateTime.is_leap_year(y))
            if d < k:
                break
            m += 1
            d -= k
        if m == 13:
            y += 1
            m = 1
        d += 1
        return DateTime(y, m, d)

    @staticmethod
    def from_ymd(s: str, sep: str = "-") -> "DateTime":
        y, m, d = s.split(sep)
        return DateTime(int(y), int(m), int(d))

    @staticmethod
    def is_leap_year(y: int) -> bool:
        """判断是否为闰年."""
        if y % 400 == 0:
            return True
        return y % 4 == 0 and y % 100 != 0

    __slots__ = ("year", "month", "day")

    def __init__(self, y: int, m: int, d: int):
        self.year = y
        self.month = m
        self.day = d

    def to_int(self) -> int:
        """基准:1年1月1日为0"""
        y = self.year - 1 if self.month <= 2 else self.year
        m = self.month + 12 if self.month <= 2 else self.month
        d = self.day
        return 365 * y + y // 4 - y // 100 + y // 400 + 306 * (m + 1) // 10 + d - 429

    def weekday(self) -> int:
        """基准:星期日为0,取值范围[0,7)"""
        return (self.to_int() + 1) % 7

    def to_string(self, sep: str = "-") -> str:
        """格式:yyyy[sep]mm[sep]dd"""
        y = str(self.year)
        m = str(self.month)
        d = str(self.day)
        while len(y) < 4:
            y = "0" + y
        while len(m) < 2:
            m = "0" + m
        while len(d) < 2:
            d = "0" + d
        return y + sep + m + sep + d

    def _to_tuple(self) -> Tuple[int, int, int]:
        return self.year, self.month, self.day

    def __iadd__(self, other: int) -> "DateTime":
        self.day += other
        lim = DateTime._MONTH_DAYS[self.month]
        if DateTime.is_leap_year(self.year) and self.month == 2:
            lim = 29
        if self.day <= lim:
            return self
        self.day = 1
        self.month += 1
        if self.month == 13:
            self.year += 1
            self.month = 1
        return self

    def __add__(self, other: int) -> "DateTime":
        return DateTime(self.year, self.month, self.day) + other

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, DateTime):
            return False
        return self._to_tuple() == other._to_tuple()

    def __ne__(self, other: object) -> bool:
        return not self == other

    def __lt__(self, other: object) -> bool:
        if not isinstance(other, DateTime):
            return False
        return self._to_tuple() < other._to_tuple()

    def __le__(self, other: object) -> bool:
        if not isinstance(other, DateTime):
            return False
        return self._to_tuple() <= other._to_tuple()

    def __gt__(self, other: object) -> bool:
        if not isinstance(other, DateTime):
            return False
        return self._to_tuple() > other._to_tuple()

    def __ge__(self, other: object) -> bool:
        if not isinstance(other, DateTime):
            return False
        return self._to_tuple() >= other._to_tuple()

    def __repr__(self) -> str:
        return f"DateTime({self.year}, {self.month}, {self.day})"


if __name__ == "__main__":
    # https://leetcode.cn/problems/reformat-date/
    # 1507. 转变日期格式
    # 请你将字符串转变为 YYYY-MM-DD 的格式，其中：
    class Solution:
        def reformatDate(self, date: str) -> str:
            day, month, year = date.split()
            month = "JanFebMarAprMayJunJulAugSepOctNovDec".index(month[:3]) // 3 + 1
            day = int(day[:-2])
            return f"{year}-{month:02d}-{day:02d}"

    # https://leetcode.cn/problems/number-of-days-between-two-dates/
    # 1360. 日期之间隔几天
    class Solution2:
        def daysBetweenDates(self, date1: str, date2: str) -> int:
            d1 = DateTime.from_ymd(date1)
            d2 = DateTime.from_ymd(date2)
            return abs(d1.to_int() - d2.to_int())
