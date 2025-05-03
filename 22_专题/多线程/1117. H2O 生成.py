# 1117. H2O 生成
# https://leetcode.cn/problems/building-h2o/description/?envType=problem-list-v2&envId=concurrency
# 输入: water = "HOH"
# 输出: "HHO"
# 解释: "HOH" 和 "OHH" 依然都是有效解。


import threading
from typing import Callable


class H2O:
    def __init__(self):
        self.barrier = threading.Barrier(3)
        self.h = threading.BoundedSemaphore(2)
        self.o = threading.Lock()  # Mutex等效于val为1的semaphore

    def hydrogen(self, releaseHydrogen: "Callable[[], None]") -> None:
        with self.h:  # 获取信号量，如果已有2个H线程在执行，第3个会被阻塞
            # releaseHydrogen() outputs "H". Do not change or remove this line.
            releaseHydrogen()
            self.barrier.wait()

    def oxygen(self, releaseOxygen: "Callable[[], None]") -> None:
        with self.o:  # 获取锁，如果已有1个O线程在执行，第2个会被阻塞
            # releaseOxygen() outputs "O". Do not change or remove this line.
            releaseOxygen()
            self.barrier.wait()
