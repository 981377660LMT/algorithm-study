from enum import Enum


class VectorClock:
    """向量时钟"""

    __slots__ = "clock", "process_id"

    def __init__(self, num_processes: int, process_id: int):
        self.clock = [0] * num_processes
        self.process_id = process_id

    def tick(self):
        """进程内部事件，时钟加1"""
        self.clock[self.process_id] += 1

    def send_event(self):
        """发送事件，时钟加1并返回当前时钟"""
        self.tick()
        return self.clock.copy()

    def receive_event(self, received_clock: list[int]):
        """接收事件，更新时钟"""
        for i in range(len(self.clock)):
            self.clock[i] = max(self.clock[i], received_clock[i])
        self.tick()

    def __str__(self):
        return str(self.clock)


class Order(Enum):
    LESS = 0
    GREATER = 1
    EQUAL = 2
    CONCURRENT = 3


def compare_clock(clock1: "VectorClock", clock2: "VectorClock") -> Order:
    """比较两个向量时钟"""
    less = False
    greater = False
    for i in range(len(clock1.clock)):
        if clock1.clock[i] < clock2.clock[i]:
            less = True
        elif clock1.clock[i] > clock2.clock[i]:
            greater = True
    if less and not greater:
        return Order.LESS
    elif greater and not less:
        return Order.GREATER
    elif not less and not greater:
        return Order.EQUAL
    else:
        return Order.CONCURRENT


# 示例使用
if __name__ == "__main__":
    num_processes = 3
    p1 = VectorClock(num_processes, 0)
    p2 = VectorClock(num_processes, 1)
    p3 = VectorClock(num_processes, 2)

    # 进程1发生内部事件
    p1.tick()
    print(f"p1: {p1}")

    # 进程1发送消息给进程2
    msg = p1.send_event()
    print(f"p1 sends message: {msg}")

    # 进程2接收消息
    p2.receive_event(msg)
    print(f"p2 receives message: {p2}")

    # 进程2发送消息给进程3
    msg = p2.send_event()
    print(f"p2 sends message: {msg}")

    # 进程3接收消息
    p3.receive_event(msg)
    print(f"p3 receives message: {p3}")

    # 进程3发生内部事件
    p3.tick()
    print(f"p3: {p3}")

    print(compare_clock(p1, p2))
    print(compare_clock(p2, p3))
    print(compare_clock(p1, p3))
