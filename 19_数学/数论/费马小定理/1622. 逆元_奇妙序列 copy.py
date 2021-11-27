MOD = 10 ** 9 + 7


class Fancy(object):
    def __init__(self):
        self.nums = []
        self.add = [0]
        self.mul = [1]

    def append(self, val: int) -> None:
        self.nums.append(val)
        self.add.append(self.add[-1])
        self.mul.append(self.mul[-1])

    def addAll(self, inc: int) -> None:
        self.add[-1] += inc

    def multAll(self, m: int) -> None:
        self.add[-1] = self.add[-1] * m % MOD
        self.mul[-1] = self.mul[-1] * m % MOD

    def getIndex(self, idx: int) -> int:
        if idx >= len(self.nums):
            return -1

        逆元 = pow(self.mul[idx], MOD - 2, MOD)
        乘的倍数 = self.mul[-1] * 逆元
        加的大小 = self.add[-1] - self.add[idx] * 乘的倍数
        return (self.nums[idx] * 乘的倍数 + 加的大小) % MOD
