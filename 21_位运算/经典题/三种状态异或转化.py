class Solution:
    def solve(self, quxes):
        n = len(quxes)
        if len(set(quxes)) == 1:
            return n
        if n <= 1:
            return n

        xor_ = 0
        mapping = {"R": 1, "G": 2, "B": 3}
        for qux in quxes:
            xor_ ^= mapping[qux]
        return 2 if xor_ == 0 else 1


print(Solution().solve(quxes=["R", "G", "B", "G", "B"]))


# 两种不同颜色在一起可以转第三种颜色
# 问最后剩下几个数
# 对应 1 2 3 的异或

# 01
# 10
# 11
