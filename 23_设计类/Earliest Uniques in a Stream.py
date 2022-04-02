# 数据流中的最早的唯一数
class EarliestUnique:
    def __init__(self, nums):
        self.visited = set()
        self.unique = dict()
        for num in nums:
            self.add(num)

    def add(self, num):
        if num in self.visited:
            return

        if num not in self.unique:
            self.unique[num] = 1
            return

        self.unique.pop(num)
        self.visited.add(num)

    def earliestUnique(self):
        if not self.unique:
            return -1

        # 注意set()不是插入有序的 dict才是插入有序的
        return next(iter(self.unique))


eu = EarliestUnique([2, 1, 3])
eu.add(3)
print(eu.__dict__)
