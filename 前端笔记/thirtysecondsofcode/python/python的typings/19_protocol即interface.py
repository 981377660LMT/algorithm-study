from abc import abstractmethod
from typing import Optional
from typing_extensions import Protocol, runtime_checkable

# 用于把 Protocol 类标记为运行时协议。
# @runtime_checkable
class TreeLike(Protocol):
    value: int

    @property
    def left(self) -> Optional['TreeLike']:
        ...

    @property
    def right(self) -> Optional['TreeLike']:
        ...

    @abstractmethod
    def walk(self) -> None:
        ...


class SimpleTree:
    def __init__(self, value: int) -> None:
        self.value = value
        self.left: Optional['SimpleTree'] = None
        self.right: Optional['SimpleTree'] = None
        self.walk = lambda: None


root: TreeLike = SimpleTree(0)  # OK
