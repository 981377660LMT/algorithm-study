from sortedcontainers import SortedList


class AllOne:
    def __init__(self):
        self.counter = dict()
        self.sl = SortedList()

    def inc(self, key: str) -> None:
        if key not in self.counter:
            self.counter[key] = 1
            self.sl.add((1, key))
        else:
            self.sl.remove((self.counter[key], key))
            self.counter[key] += 1
            self.sl.add((self.counter[key], key))

    def dec(self, key: str) -> None:
        self.sl.remove((self.counter[key], key))
        self.counter[key] -= 1
        if self.counter[key] == 0:
            del self.counter[key]
        else:
            self.sl.add((self.counter[key], key))

    def getMaxKey(self) -> str:
        if not self.sl:
            return ""
        return self.sl[-1][1]

    def getMinKey(self) -> str:
        if not self.sl:
            return ""
        return self.sl[0][1]


# Your AllOne object will be instantiated and called as such:
# obj = AllOne()
# obj.inc(key)
# obj.dec(key)
# param_3 = obj.getMaxKey()
# param_4 = obj.getMinKey()
