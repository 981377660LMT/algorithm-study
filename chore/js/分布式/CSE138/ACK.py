from dataclasses import dataclass
from typing import Optional
import time
import threading
import queue


@dataclass
class Message:
    sender_id: str
    sequence_number: int
    content: str


@dataclass
class Ack:
    sender_id: str
    sequence_number: int


class FIFOChannel:
    """模拟 FIFO 通道，传输消息和确认"""

    def __init__(self):
        self.sender_to_receiver = queue.Queue()
        self.receiver_to_sender = queue.Queue()

    def transmit_message(self, message: Message):
        """传输消息到接收者"""
        print(f"Channel: Transmitting message {message}")
        self.sender_to_receiver.put(message)

    def transmit_ack(self, ack: Ack):
        """传输确认到发送者"""
        print(f"Channel: Transmitting ACK {ack}")
        self.receiver_to_sender.put(ack)

    def receive_message(self, timeout: Optional[float] = None) -> Optional[Message]:
        """接收来自发送者的消息"""
        try:
            message = self.sender_to_receiver.get(timeout=timeout)
            return message
        except queue.Empty:
            return None

    def receive_ack(self, timeout: Optional[float] = None) -> Optional[Ack]:
        """接收来自接收者的确认"""
        try:
            ack = self.receiver_to_sender.get(timeout=timeout)
            return ack
        except queue.Empty:
            return None


class Sender:
    """发送者，维护序列号并发送消息，等待确认"""

    def __init__(self, sender_id: str, channel: FIFOChannel, timeout: float = 2.0):
        self.sender_id = sender_id
        self.channel = channel
        self.sequence_number = 0
        self.timeout = timeout
        self.lock = threading.Lock()
        self.ack_event = threading.Event()
        self.expected_ack = -1

        # 启动接收ACK线程
        self.stop_event = threading.Event()
        self.ack_thread = threading.Thread(target=self._receive_acks, daemon=True)
        self.ack_thread.start()

    def send(self, content: str):
        """发送消息并等待确认"""
        with self.lock:
            message = Message(
                sender_id=self.sender_id, sequence_number=self.sequence_number, content=content
            )
            self.expected_ack = self.sequence_number
            self.ack_event.clear()
            self.channel.transmit_message(message)
            print(f"Sender {self.sender_id}: Sent {message}")

            # 等待ACK
            ack_received = self.ack_event.wait(timeout=self.timeout)
            if ack_received:
                print(f"Sender {self.sender_id}: ACK received for seq {self.sequence_number}")
                self.sequence_number += 1
            else:
                print(
                    f"Sender {self.sender_id}: ACK timeout for seq {self.sequence_number}, resending..."
                )
                # 重传消息
                self.channel.transmit_message(message)
                print(f"Sender {self.sender_id}: Resent {message}")
                # 等待再次ACK
                ack_received = self.ack_event.wait(timeout=self.timeout)
                if ack_received:
                    print(f"Sender {self.sender_id}: ACK received for seq {self.sequence_number}")
                    self.sequence_number += 1
                else:
                    print(
                        f"Sender {self.sender_id}: Failed to receive ACK for seq {self.sequence_number}"
                    )

    def stop(self):
        """停止接收ACK线程"""
        self.stop_event.set()
        self.ack_thread.join()

    def _receive_acks(self):
        """接收ACK并设置事件"""
        while not self.stop_event.is_set():
            ack = self.channel.receive_ack(timeout=1.0)
            if ack and ack.sender_id == self.sender_id:
                with self.lock:
                    if ack.sequence_number == self.expected_ack:
                        self.ack_event.set()


class Receiver:
    """接收者，维护预期序列号并发送确认"""

    def __init__(self, channel: FIFOChannel):
        self.channel = channel
        self.expected_seq: int = 0
        self.lock = threading.Lock()

        # 启动接收消息线程
        self.stop_event = threading.Event()
        self.receive_thread = threading.Thread(target=self._receive_messages, daemon=True)
        self.receive_thread.start()

    def _receive_messages(self):
        """接收并处理消息"""
        while not self.stop_event.is_set():
            message = self.channel.receive_message(timeout=1.0)
            if message:
                with self.lock:
                    print(f"Receiver: Received {message}")
                    if message.sequence_number == self.expected_seq:
                        self._deliver(message)
                        ack = Ack(
                            sender_id=message.sender_id, sequence_number=message.sequence_number
                        )
                        self.channel.transmit_ack(ack)
                        print(f"Receiver: Sent ACK for seq {ack.sequence_number}")
                        self.expected_seq += 1
                    elif message.sequence_number < self.expected_seq:
                        # 重复或过时的消息，发送ACK
                        ack = Ack(
                            sender_id=message.sender_id, sequence_number=message.sequence_number
                        )
                        self.channel.transmit_ack(ack)
                        print(f"Receiver: Resent ACK for seq {ack.sequence_number}")
                    else:
                        # 未来的消息，暂时不处理或请求重传
                        print(
                            f"Receiver: Out-of-order message {message}, expected seq {self.expected_seq}. Ignoring."
                        )

    def stop(self):
        """停止接收消息线程"""
        self.stop_event.set()
        self.receive_thread.join()

    def _deliver(self, message: Message):
        """交付消息到应用层"""
        print(f"Receiver: Delivered {message}")


# 示例使用
if __name__ == "__main__":
    channel = FIFOChannel()
    receiver = Receiver(channel)
    sender = Sender("P1", channel, timeout=3.0)

    try:
        # 发送顺序消息
        sender.send("Message 1")
        sender.send("Message 2")
        sender.send("Message 3")

        print("\n--- 模拟ACK丢失导致重传 ---\n")

        # 修改FIFOChannel以模拟ACK丢失
        original_transmit_ack = channel.transmit_ack

        def unreliable_transmit_ack(ack: Ack):
            if ack.sequence_number == 1:
                # 模拟ACK丢失
                print(f"Channel: ACK {ack} lost!")
                return
            original_transmit_ack(ack)

        channel.transmit_ack = unreliable_transmit_ack

        sender.send("Message 4")  # ACK丢失，触发重传
        sender.send("Message 5")
    finally:
        # 等待一段时间以处理消息
        time.sleep(5)

        # 停止线程
        sender.stop()
        receiver.stop()
