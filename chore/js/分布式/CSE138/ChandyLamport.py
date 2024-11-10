import threading
import time
from typing import Callable, Dict, List


class Process:
    def __init__(self, id: int, outgoing: List[int], send_marker: Callable[[int], None]):
        self.id = id
        self.outgoing = outgoing
        self.send_marker = send_marker
        self.state = f"State of process {id}"
        self.recorded_state = None
        self.channel_states: Dict[int, List[str]] = {ch: [] for ch in outgoing}
        self.has_recorded = False
        self.lock = threading.Lock()

    def receive_message(self, from_process: int, message: str):
        with self.lock:
            if message == "MARKER":
                self.handle_marker(from_process)
            else:
                if not self.has_recorded:
                    self.channel_states[from_process].append(message)
                print(f"Process {self.id} received message from Process {from_process}: {message}")

    def handle_marker(self, from_process: int):
        if not self.has_recorded:
            self.recorded_state = self.state
            self.has_recorded = True
            print(f"Process {self.id} records its state: {self.recorded_state}")
            for ch in self.outgoing:
                self.send_marker(ch)
        # Mark the channel as empty
        self.channel_states[from_process] = []
        print(f"Process {self.id} marks channel from Process {from_process} as empty")

    def send_message(self, to_process: int, message: str):
        print(f"Process {self.id} sends message to Process {to_process}: {message}")
        processes[to_process].receive_message(self.id, message)

    def initiate_snapshot(self):
        print(f"Process {self.id} initiates snapshot")
        self.handle_marker(self.id)


def send_marker(to_process: int):
    processes[to_process].receive_message(-1, "MARKER")


# 创建进程
processes: Dict[int, Process] = {}
num_processes = 3
connections = {0: [1, 2], 1: [0, 2], 2: [0, 1]}

for pid in range(num_processes):
    processes[pid] = Process(pid, connections[pid], send_marker)


# 模拟消息传递
def simulate():
    processes[0].send_message(1, "Message 1 from P0 to P1")
    processes[1].send_message(2, "Message 1 from P1 to P2")
    time.sleep(1)
    processes[0].initiate_snapshot()
    processes[2].send_message(0, "Message 1 from P2 to P0")
    time.sleep(1)


thread = threading.Thread(target=simulate)
thread.start()
thread.join()
