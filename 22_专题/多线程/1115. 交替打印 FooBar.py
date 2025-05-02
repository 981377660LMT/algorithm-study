import threading
from typing import Callable


class FooBar:
    def __init__(self, n):
        self.n = n
        self.l1, self.l2 = threading.Lock(), threading.Lock()

        self.l2.acquire()

    def foo(self, printFoo: "Callable[[], None]") -> None:
        for _ in range(self.n):
            self.l1.acquire()
            printFoo()
            self.l2.release()

    def bar(self, printBar: "Callable[[], None]") -> None:
        for _ in range(self.n):
            self.l2.acquire()
            printBar()
            self.l1.release()
