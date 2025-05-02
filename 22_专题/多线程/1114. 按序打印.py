from threading import Event
from typing import Callable


class Foo:
    def __init__(self):
        self.first_done = Event()
        self.second_done = Event()

    def first(self, printFirst: "Callable[[], None]") -> None:

        # printFirst() outputs "first". Do not change or remove this line.
        printFirst()
        self.first_done.set()

    def second(self, printSecond: "Callable[[], None]") -> None:
        self.first_done.wait()
        # printSecond() outputs "second". Do not change or remove this line.
        printSecond()
        self.second_done.set()

    def third(self, printThird: "Callable[[], None]") -> None:
        self.second_done.wait()
        # printThird() outputs "third". Do not change or remove this line.
        printThird()
