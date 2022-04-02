from collections import defaultdict
from sortedcontainers import SortedList


class HistoricalMap:
    def __init__(self):
        self.record = defaultdict(lambda: SortedList())

    def get(self, key, timestamp):
        if key not in self.record:
            return -1
        pos = self.record[key].bisect_right((timestamp, int(1e20))) - 1
        if pos < 0:
            return -1
        return self.record[key][pos][1]

    def set(self, key, val, timestamp):
        # set前要看是否修改之前的记录
        pos = self.record[key].bisect_right((timestamp, int(1e20))) - 1
        if pos >= 0 and self.record[key][pos][0] == timestamp:
            self.record[key].pop(pos)
        self.record[key].add((timestamp, val))


if __name__ == '__main__':
    hMap = HistoricalMap()
    hMap.set(1, 3, 1)
    hMap.set(1, 2, 1)
    print(hMap.get(1, 5))

