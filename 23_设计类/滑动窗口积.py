# 滑动窗口后k个数的乘积
class SlidingWindowProduct:
    def __init__(self):
        self.products = [1]

    def add(self, num):
        if num:
            self.products.append(self.products[-1] * num)
        else:
            self.products = [1]

    def product(self, k):
        if len(self.products) > k:
            return self.products[-1] // self.products[-k - 1]
        return 0
