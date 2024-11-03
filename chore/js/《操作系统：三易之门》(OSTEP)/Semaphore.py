import threading
import time
import random

# 创建一个信号量，最多允许3个线程同时访问
semaphore = threading.Semaphore(3)

# 共享资源
shared_resource = []


def worker(thread_id):
    print(f"线程 {thread_id} 正在等待获取信号量...")
    semaphore.acquire()
    try:
        print(f"线程 {thread_id} 获取了信号量，开始工作。")
        # 模拟工作时间
        time.sleep(random.uniform(1, 3))
        shared_resource.append(thread_id)
        print(f"线程 {thread_id} 完成工作。")
    finally:
        semaphore.release()
        print(f"线程 {thread_id} 释放了信号量。")


if __name__ == "__main__":
    threads = []
    for i in range(10):
        t = threading.Thread(target=worker, args=(i,))
        threads.append(t)
        t.start()

    for t in threads:
        t.join()

    print("所有线程已完成。")
    print("共享资源内容:", shared_resource)
