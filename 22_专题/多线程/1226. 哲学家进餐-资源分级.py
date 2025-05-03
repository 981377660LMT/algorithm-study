# !操作系统哲学家吃饭的几种解法
#
# - 限制进餐人数(破除环)：最多允许4个哲学家同时尝试拿叉子
# - 资源分级：为每个叉子分配一个编号，哲学家总是先拿起编号较小的叉子
# - 奇偶解法：奇数哲学家先拿左边的叉子，偶数哲学家先拿右边的叉子
# - 串行化


import threading
from typing import Callable


class DiningPhilosophers:

    def __init__(self):
        self.locks = [threading.Lock() for _ in range(5)]

    # call the functions directly to execute, for example, eat()
    def wantsToEat(
        self,
        philosopher: int,
        pickLeftFork: "Callable[[], None]",
        pickRightFork: "Callable[[], None]",
        eat: "Callable[[], None]",
        putLeftFork: "Callable[[], None]",
        putRightFork: "Callable[[], None]",
    ) -> None:
        if philosopher == 0:
            right = philosopher
            left = (philosopher + 1) % 5
        else:
            left = philosopher
            right = (philosopher + 1) % 5

        with self.locks[left], self.locks[right]:
            pickLeftFork()
            pickRightFork()
            eat()
            putLeftFork()
            putRightFork()
