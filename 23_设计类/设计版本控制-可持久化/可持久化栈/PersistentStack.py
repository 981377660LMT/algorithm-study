# push(x): データxをスタックの一番上に格納した後のスタックを返す（このスタックは変化なし）
# top(): スタックの一番上にあるデータを返す
# pop(): スタックの一番上のデータを削除した後のスタックを返す（このスタックは変化なし）

# 永続スタックの実装方針
# 連結リストを用いて全バージョンのデータを管理します。
# 連結リストの各ノードには二つのデータ(value, pre)を保持させます。


from typing import Generic, Optional, TypeVar


E = TypeVar("E")


class PersistentStack(Generic[E]):
    """fully persistent stack"""

    __slots__ = ("value", "pre", "_index")

    @staticmethod
    def default() -> "PersistentStack[E]":
        """return an empty stack whose pre is itself"""
        res = PersistentStack(None, None, 0)
        res.pre = res
        return res

    def __init__(
        self, value: Optional["E"] = None, pre: Optional["PersistentStack[E]"] = None, index=0
    ) -> None:
        self.value = value
        self.pre = pre
        self._index = index

    def top(self) -> Optional["E"]:
        """return the top element"""
        return self.value

    def push(self, x: "E") -> "PersistentStack[E]":
        """push x to the top of the stack and return a new stack"""
        return PersistentStack(x, self, self._index + 1)

    def pop(self) -> "PersistentStack[E]":
        """pop the top element and return previous stack"""
        if self.pre is None:
            raise IndexError("pop from empty stack")
        return self.pre

    def __repr__(self) -> str:
        return f"PersistentStack(index={self._index}, value={self.value})"

    def __len__(self) -> int:
        return self._index


if __name__ == "__main__":

    stack0 = PersistentStack()
    stack0.pre = stack0  # !trick:可持久化栈的栈底元素指向自身
    stack1 = stack0.push(1)
    print(stack1)  # 1
    stack2 = stack1.push(2)
    print(stack2)  # 2
    print(stack0.pop())
