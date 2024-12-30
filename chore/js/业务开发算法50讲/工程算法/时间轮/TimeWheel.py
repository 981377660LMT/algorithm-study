import time
import threading


class TimerTask:
    def __init__(self, timeout, callback):
        self.expire = timeout
        self.callback = callback


class TimeWheel:
    def __init__(self, tick, wheel_size):
        self.tick = tick
        self.wheel_size = wheel_size
        self.current_slot = 0
        self.slots = [[] for _ in range(wheel_size)]
        self.lock = threading.Lock()
        self.running = False

    def add_task(self, task):
        with self.lock:
            ticks = task.expire // self.tick
            slot = (self.current_slot + ticks) % self.wheel_size
            self.slots[slot].append(task)

    def tick_handler(self):
        while self.running:
            time.sleep(self.tick)
            with self.lock:
                tasks = self.slots[self.current_slot]
                self.slots[self.current_slot] = []
            for task in tasks:
                task.callback()
            self.current_slot = (self.current_slot + 1) % self.wheel_size

    def start(self):
        self.running = True
        threading.Thread(target=self.tick_handler, daemon=True).start()

    def stop(self):
        self.running = False


# 使用示例
def my_task():
    print("Task executed at", time.time())


if __name__ == "__main__":
    tw = TimeWheel(tick=1, wheel_size=60)
    tw.start()
    tw.add_task(TimerTask(timeout=5, callback=my_task))
    tw.add_task(TimerTask(timeout=10, callback=my_task))
    time.sleep(15)
    tw.stop()
