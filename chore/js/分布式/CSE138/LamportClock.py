class LamportClock:
    """Lamport 时钟"""

    __slots__ = "clock"

    def __init__(self):
        self.clock = 0

    def tick(self):
        """进程内部事件，时钟加1"""
        self.clock += 1

    def send_event(self):
        """发送事件，时钟加1并返回当前时钟"""
        self.tick()
        return self.clock

    def receive_event(self, received_clock: int):
        """接收事件，更新时钟"""
        self.clock = max(self.clock, received_clock) + 1

    def __str__(self):
        return str(self.clock)


# 示例使用
if __name__ == "__main__":
    p1 = LamportClock()
    p2 = LamportClock()

    # 进程1发生内部事件
    p1.tick()
    print(f"p1: {p1}")

    # 进程1发送消息给进程2
    msg = p1.send_event()
    print(f"p1 sends message: {msg}")

    # 进程2接收消息
    p2.receive_event(msg)
    print(f"p2 receives message: {p2}")

    # 进程2发生内部事件
    p2.tick()
    print(f"p2: {p2}")

    # 进程2发送消息给进程1
    msg = p2.send_event()
    print(f"p2 sends message: {msg}")

    # 进程1接收消息
    p1.receive_event(msg)
    print(f"p1 receives message: {p1}")
