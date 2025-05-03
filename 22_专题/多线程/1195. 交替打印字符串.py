# FizzBuzz 问题要求：
#
# 打印从 1 到 n 的数字
# 如果数字能被 3 整除，打印 "Fizz"
# 如果数字能被 5 整除，打印 "Buzz"
# 如果数字同时能被 3 和 5 整除，打印 "FizzBuzz"
# 其他情况，直接打印数字
# 这个多线程版本通过四个不同的线程来完成这些任务。
#
# !核心思想是：用 number 线程作为控制者，它先获取锁，然后根据当前数字决定释放哪个线程的信号量。

from typing import Callable
from threading import Semaphore


class FizzBuzz:
    def __init__(self, n: int):
        self.n = n + 1
        self.lock_X = Semaphore(1)  # number 线程的信号量，初始值为1表示开始时可以执行
        self.lock_F = Semaphore(0)
        self.lock_B = Semaphore(0)
        self.lock_FB = Semaphore(0)

    def number(self, printNumber: "Callable[[int], None]") -> None:
        for i in range(1, self.n):
            self.lock_X.acquire()  # 获取控制权
            if i % 3 == i % 5 == 0:  # 如果能同时被3和5整除
                self.lock_FB.release()  # 释放 FizzBuzz 线程
            elif i % 3 == 0:
                self.lock_F.release()
            elif i % 5 == 0:
                self.lock_B.release()
            else:
                printNumber(i)
                self.lock_X.release()

    def fizz(self, printFizz: "Callable[[], None]") -> None:
        for i in range(3, self.n, 3):
            if i % 5:
                self.lock_F.acquire()  # 等待 number 线程的信号
                printFizz()
                self.lock_X.release()  # 释放控制权给 number 线程

    def buzz(self, printBuzz: "Callable[[], None]") -> None:
        for i in range(5, self.n, 5):
            if i % 3:
                self.lock_B.acquire()
                printBuzz()
                self.lock_X.release()

    def fizzbuzz(self, printFizzBuzz: "Callable[[], None]") -> None:
        for _ in range(15, self.n, 15):
            self.lock_FB.acquire()
            printFizzBuzz()
            self.lock_X.release()


if __name__ == "__main__":
    import threading

    n = 15
    fizzbuzz = FizzBuzz(n)

    def printFizz():
        print("fizz", end=" ")

    def printBuzz():
        print("buzz", end=" ")

    def printFizzBuzz():
        print("fizzbuzz", end=" ")

    def printNumber(x):
        print(x, end=" ")

    threads = [
        threading.Thread(target=fizzbuzz.fizz, args=(printFizz,)),
        threading.Thread(target=fizzbuzz.buzz, args=(printBuzz,)),
        threading.Thread(target=fizzbuzz.fizzbuzz, args=(printFizzBuzz,)),
        threading.Thread(target=fizzbuzz.number, args=(printNumber,)),
    ]

    # Start all threads
    for thread in threads:
        thread.start()

    # Wait for all threads to finish
    for thread in threads:
        thread.join()
