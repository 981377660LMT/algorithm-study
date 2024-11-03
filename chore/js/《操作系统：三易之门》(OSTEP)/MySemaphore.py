import threading


class Semaphore:
    def __init__(self, value=1):
        self.value = value
        self.lock = threading.Lock()
        self.condition = threading.Condition(self.lock)

    def acquire(self, blocking=True, timeout=None):
        with self.condition:
            if not blocking and self.value <= 0:
                return False
            end_time = None
            if timeout is not None:
                end_time = time.time() + timeout
            while self.value <= 0:
                if not blocking:
                    return False
                remaining = None
                if end_time is not None:
                    remaining = end_time - time.time()
                    if remaining <= 0:
                        return False
                if not self.condition.wait(timeout=remaining):
                    return False
            self.value -= 1
            return True

    def release(self):
        with self.condition:
            self.value += 1
            self.condition.notify()


if __name__ == "__main__":

    # 示例使用
    import time
    import random

    def worker(semaphore, thread_id):
        print(f"线程 {thread_id} 尝试获取信号量...")
        acquired = semaphore.acquire(timeout=2)
        if acquired:
            print(f"线程 {thread_id} 获取到信号量，开始工作。")
            time.sleep(random.uniform(1, 3))
            semaphore.release()
            print(f"线程 {thread_id} 释放了信号量。")
        else:
            print(f"线程 {thread_id} 未能在2秒内获取到信号量。")

    semaphore = Semaphore(3)  # 最多允许3个线程同时访问
    threads = []
    for i in range(5):
        t = threading.Thread(target=worker, args=(semaphore, i))
        threads.append(t)
        t.start()

    for t in threads:
        t.join()

    print("所有线程已完成。")
