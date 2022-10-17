MOD = int(1e9 + 7)

# 奇妙序列


class Fancy(object):
    def __init__(self):
        self.nums = []
        self.add = [0]
        self.mul = [1]

    def append(self, val: int) -> None:
        """将 val 添加到 Fancy 序列末尾"""
        self.nums.append(val)
        self.add.append(self.add[-1])
        self.mul.append(self.mul[-1])

    def addAll(self, inc: int) -> None:
        """将 inc 添加到每个数上"""
        self.add[-1] += inc

    def multAll(self, m: int) -> None:
        """将每个数乘以 m"""
        self.add[-1] = self.add[-1] * m % MOD
        self.mul[-1] = self.mul[-1] * m % MOD

    def getIndex(self, idx: int) -> int:
        """返回下标为 idx 的数 模 10^9 + 7
        如果下标为 idx 的数不存在，返回 -1
        """
        if idx >= len(self.nums):
            return -1

        inv = pow(self.mul[idx], MOD - 2, MOD)
        res1 = self.mul[-1] * inv
        res2 = self.add[-1] - self.add[idx] * res1
        return (self.nums[idx] * res1 + res2) % MOD
