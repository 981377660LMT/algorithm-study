from collections import defaultdict


class BooleanArray:
    def __init__(self):
        self.arr = defaultdict(lambda: False)

    def setTrue(self, i):
        self.arr[i] = True

    def setFalse(self, i):
        self.arr[i] = False

    def setAllTrue(self):
        self.arr = defaultdict(lambda: True)

    def setAllFalse(self):
        self.arr = defaultdict(lambda: False)

    def getValue(self, i):
        return self.arr[i]
