from math import exp
from random import shuffle, randint, random
from typing import List

# TLE


class Solution:
    def maxHappyGroups(self, batchSize: int, groups: List[int]) -> int:
        w = []  # 数组
        n = 0  # 数组长度
        iter_epoch = 40  # 模拟退火迭代次数
        anneal_rate = 0.985  # 退火率
        self.res = 0  # 数组计算出的结果

        def calc() -> int:
            happy_group = 1
            cur_sum = 0
            for i in range(n):
                cur_sum = (cur_sum + w[i]) % batchSize
                if cur_sum == 0 and i != n - 1:
                    happy_group += 1
            self.res = max(self.res, happy_group)
            return happy_group

        def simulate_anneal() -> None:
            shuffle(w)
            t = 1e6
            while t > 1e-5:
                a = randint(0, 2 ** 10) % n
                b = randint(0, 2 ** 10) % n
                x = calc()
                w[a], w[b] = w[b], w[a]
                y = calc()
                delta = y - x
                if delta >= 0:
                    continue
                else:
                    if exp(delta / t) > random():
                        w[a], w[b] = w[b], w[a]
                # ---- 退火
                t *= anneal_rate

        if batchSize == 1:
            return len(groups)

        already_ok = 0  # 这个组是个模块。直接放到最前面
        for g in groups:
            if g % batchSize == 0:
                already_ok += 1
            else:
                w.append(g % batchSize)

        n = len(w)

        # -- 就1组或者0组，直接返回
        if n <= 1:
            return already_ok + n

        for _ in range(iter_epoch):
            simulate_anneal()

        return self.res + already_ok

