from typing import Generic, Optional, TypeVar

V = TypeVar("V")


class Dictionary(Generic[V]):
    """获取对象唯一标识的字典."""

    __slots__ = "_valueToId", "_idToValue"

    def __init__(self):
        self._valueToId = dict()
        self._idToValue = []

    def id(self, value: V) -> int:
        res = self._valueToId.get(value, None)
        if res is not None:
            return res
        id_ = len(self._idToValue)
        self._idToValue.append(value)
        self._valueToId[value] = id_
        return id_

    def value(self, id_: int) -> Optional[V]:
        if id_ < 0 or id_ >= len(self._idToValue):
            return None
        return self._idToValue[id_]

    def __len__(self) -> int:
        return len(self._idToValue)


if __name__ == "__main__":
    d = Dictionary[str]()
    print(d.id("a"))
    print(d.id("b"))
    print(d.id("a"))
    print(d.value(0))
    print(d.value(1))
    print(d.value(2))
    print(len(d))
