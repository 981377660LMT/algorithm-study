from typing import Dict
from collections import defaultdict
from dataclasses import dataclass

from sortedcontainers import SortedList


@dataclass
class Message:
    sender_id: str
    sequence_number: int
    content: str


class Sender:
    """发送者，维护序列号并发送消息"""

    __slots__ = "sender_id", "sequence_number", "channel"

    def __init__(self, sender_id: str, channel: "FIFOChannel"):
        self.sender_id = sender_id
        self.sequence_number = 0
        self.channel = channel

    def send(self, content: str):
        """发送消息，序列号递增"""
        message = Message(
            sender_id=self.sender_id, sequence_number=self.sequence_number, content=content
        )
        print(f"Sender {self.sender_id} sends: {message}")
        self.channel.transmit(message)
        self.sequence_number += 1


class Receiver:
    """接收者，维护每个发送者的预期序列号和缓冲区"""

    __slots__ = "expected_seq", "buffer"

    def __init__(self):
        self.expected_seq: Dict[str, int] = defaultdict(int)
        self.buffer: Dict[str, SortedList] = defaultdict(
            lambda: SortedList(key=lambda m: m.sequence_number)
        )

    def receive(self, message: Message):
        """接收并处理消息"""
        sender = message.sender_id
        seq = message.sequence_number
        expected = self.expected_seq[sender]

        print(f"Receiver got: {message}, expected_seq: {expected}")

        if seq == expected:
            self._deliver(message)
            self.expected_seq[sender] += 1

            # 检查缓冲区中是否有后续消息可以交付
            while self.buffer[sender]:
                next_msg = self.buffer[sender][0]
                if next_msg.sequence_number == self.expected_seq[sender]:
                    self.buffer[sender].pop(0)
                    self._deliver(next_msg)
                    self.expected_seq[sender] += 1
                else:
                    break
        elif seq > expected:
            # 将消息缓冲起来
            print(f"Buffering message from {sender} with seq {seq}")
            self.buffer[sender].add(message)
        else:
            # 旧消息，可能是重复的，选择丢弃
            print(f"Dropping out-of-order message from {sender} with seq {seq}")

    def _deliver(self, message: Message):
        """交付消息"""
        print(f"Delivered to application: {message}")


class FIFOChannel:
    """模拟 FIFO 通道，确保单一发送者的消息顺序"""

    def __init__(self, receiver: Receiver):
        self.receiver = receiver

    def transmit(self, message: Message):
        """传输消息到接收者，模拟可靠传输"""
        # 在这里，可以添加延迟、丢失等模拟不可靠因素
        # 为简单起见，假设传输是可靠且顺序的
        self.receiver.receive(message)


# 示例使用
if __name__ == "__main__":
    receiver = Receiver()
    channel = FIFOChannel(receiver)

    sender1 = Sender("P1", channel)
    sender2 = Sender("P2", channel)

    # 发送顺序消息
    sender1.send("m1 from P1")
    sender1.send("m2 from P1")
    sender2.send("m1 from P2")
    sender1.send("m3 from P1")
    sender2.send("m2 from P2")

    print("\n--- 模拟乱序消息传输 ---\n")

    # 修改 FIFOChannel 以模拟乱序传输
    class UnreliableFIFOChannel(FIFOChannel):
        def transmit(self, message: Message):
            """传输消息到接收者，模拟乱序"""
            if message.sequence_number == 1 and message.sender_id == "P1":
                # 故意延迟发送 m2 from P1
                print(f"Simulating delay for message: {message}")
                # 先传输 m3
                if sender1.sequence_number > 2:
                    delayed_msg = Message(
                        sender_id=message.sender_id, sequence_number=1, content=message.content
                    )
                    super().transmit(delayed_msg)
            super().transmit(message)

    # 使用不可靠的通道
    unreliable_channel = UnreliableFIFOChannel(receiver)

    sender1_unreliable = Sender("P1", unreliable_channel)
    sender2_unreliable = Sender("P2", unreliable_channel)

    sender1_unreliable.send("m1 from P1")
    sender1_unreliable.send("m2 from P1")
    sender1_unreliable.send("m3 from P1")
    sender2_unreliable.send("m1 from P2")
    sender2_unreliable.send("m2 from P2")

    # 在真实场景中，如果某条消息被延迟或丢失，接收者需要处理缓冲区
