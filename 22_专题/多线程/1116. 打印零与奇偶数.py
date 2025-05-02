import threading


class ZeroEvenOdd:
    def __init__(self, n):
        self.n = n
        # 初始化三个信号量
        self.zero_sem = threading.Semaphore(1)  # zero线程可以先执行
        self.odd_sem = threading.Semaphore(0)  # odd线程需要等待
        self.even_sem = threading.Semaphore(0)  # even线程需要等待

    def zero(self, printNumber):
        for i in range(1, self.n + 1):
            self.zero_sem.acquire()  # 获取打印0的权限
            printNumber(0)

            # 根据i的奇偶性，决定释放哪个线程
            if i % 2 == 1:
                self.odd_sem.release()
            else:
                self.even_sem.release()

    def even(self, printNumber):
        for i in range(2, self.n + 1, 2):
            self.even_sem.acquire()
            printNumber(i)
            self.zero_sem.release()

    def odd(self, printNumber):
        for i in range(1, self.n + 1, 2):
            self.odd_sem.acquire()
            printNumber(i)
            self.zero_sem.release()


# 测试代码
def print_number(x):
    print(x, end="")


def main():
    n = 5
    zeo = ZeroEvenOdd(n)

    # 创建三个线程
    t1 = threading.Thread(target=zeo.zero, args=(print_number,))
    t2 = threading.Thread(target=zeo.even, args=(print_number,))
    t3 = threading.Thread(target=zeo.odd, args=(print_number,))

    # 启动线程
    t1.start()
    t2.start()
    t3.start()

    # 等待所有线程完成
    t1.join()
    t2.join()
    t3.join()


if __name__ == "__main__":
    main()
    print()  # 打印换行
