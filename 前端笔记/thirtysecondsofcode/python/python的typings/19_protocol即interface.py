from typing import Optional
from typing_extensions import Protocol


class TreeLike(Protocol):
    value: int

    @property
    def left(self) -> Optional['TreeLike']:
        ...

    @property
    def right(self) -> Optional['TreeLike']:
        ...


class SimpleTree:
    def __init__(self, value: int) -> None:
        self.value = value
        self.left: Optional['SimpleTree'] = None
        self.right: Optional['SimpleTree'] = None


root: TreeLike = SimpleTree(0)  # OK
