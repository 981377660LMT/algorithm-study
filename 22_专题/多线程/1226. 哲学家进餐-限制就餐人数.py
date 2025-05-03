# !操作系统哲学家吃饭的几种解法
#
# - 限制进餐人数(破除环)：最多允许4个哲学家同时尝试拿叉子
# - 资源分级：为每个叉子分配一个编号，哲学家总是先拿起编号较小的叉子
# - 奇偶解法：奇数哲学家先拿左边的叉子，偶数哲学家先拿右边的叉子
# - 串行化


from typing import Callable

from threading import Lock, Semaphore


class DiningPhilosophers:
    def __init__(self):
        self.limit = Semaphore(4)  # 限制最多4个人就餐
        self.locks = [Lock() for _ in range(5)]  # 叉子锁

    def wantsToEat(
        self,
        philosopher: int,
        pickLeftFork: "Callable[[], None]",
        pickRightFork: "Callable[[], None]",
        eat: "Callable[[], None]",
        putLeftFork: "Callable[[], None]",
        putRightFork: "Callable[[], None]",
    ) -> None:
        right_fork = philosopher
        left_fork = (philosopher + 1) % 5
        with self.limit, self.locks[left_fork], self.locks[right_fork]:
            pickLeftFork()
            pickRightFork()
            eat()
            putLeftFork()
            putRightFork()
