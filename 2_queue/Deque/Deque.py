from typing import Optional, TypeVar, Iterable, Generic, List

T = TypeVar("T")


class Deque(Generic[T]):
    """
    双栈实现的可索引 deque.
    _dat1: 前半部分逆序存, _dat2: 后半部分正序.
    均摊 O(1).
    """

    __slots__ = ("_dat1", "_dat2")

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self._dat1: List[T] = []
        self._dat2: List[T] = []
        if iterable:
            self._dat2.extend(iterable)

    def append(self, x: T):
        self._dat2.append(x)

    def appendleft(self, x: T):
        self._dat1.append(x)

    def pop(self) -> T:
        if not self:
            raise IndexError("pop from an empty deque")
        if not self._dat2:
            self._rebalance_from_left()
        return self._dat2.pop()

    def popleft(self) -> T:
        if not self:
            raise IndexError("popleft from an empty deque")
        if not self._dat1:
            self._rebalance_from_right()
        return self._dat1.pop()

    def front(self) -> T:
        if not self:
            raise IndexError("front from an empty deque")
        return self[0]

    def back(self) -> T:
        if not self:
            raise IndexError("back from an empty deque")
        return self[-1]

    def clear(self):
        self._dat1.clear()
        self._dat2.clear()

    def reverse(self):
        self._dat1, self._dat2 = self._dat2, self._dat1

    def _rebalance_from_right(self):
        m = len(self._dat2)
        if m == 0:
            return
        k = (m + 1) // 2
        self._dat1 = self._dat2[k - 1 :: -1]
        self._dat2[:k] = []

    def _rebalance_from_left(self):
        m = len(self._dat1)
        if m == 0:
            return
        k = (m + 1) // 2
        self._dat2 = self._dat1[k - 1 :: -1]
        self._dat1[:k] = []

    def _norm_index(self, i: int) -> int:
        n = len(self)
        if i < 0:
            i += n
        if not (0 <= i < n):
            raise IndexError("deque index out of range")
        return i

    def __getitem__(self, i: int) -> T:
        i = self._norm_index(i)
        len1 = len(self._dat1)
        if i < len1:
            return self._dat1[len1 - 1 - i]
        return self._dat2[i - len1]

    def __setitem__(self, i: int, value: T):
        i = self._norm_index(i)
        len1 = len(self._dat1)
        if i < len1:
            self._dat1[len1 - 1 - i] = value
        else:
            self._dat2[i - len1] = value

    def __iter__(self):
        for x in reversed(self._dat1):
            yield x
        for x in self._dat2:
            yield x

    def __len__(self) -> int:
        return len(self._dat1) + len(self._dat2)

    def __bool__(self) -> bool:
        return len(self) > 0


# ==================== 对拍测试代码 ====================


def run_tests():
    import collections
    import random

    """
    通过与 collections.deque 对拍来测试 Deque 的正确性。
    """
    my_dq = Deque()
    std_dq = collections.deque()

    NUM_OPERATIONS = 2000000  # 测试操作的总次数
    MAX_VALUE = 10000  # 随机值的范围

    print(f"开始进行 {NUM_OPERATIONS} 次随机操作对拍测试...")

    for i in range(NUM_OPERATIONS):
        op_choice = random.random()

        # 检查状态是否一致
        def check_consistency(op_name: str):
            assert len(my_dq) == len(
                std_dq
            ), f"操作 {op_name} 后长度不一致: {len(my_dq)} vs {len(std_dq)}"
            assert list(my_dq) == list(std_dq), f"操作 {op_name} 后内容不一致"

        if op_choice < 0.25:  # append
            val = random.randint(0, MAX_VALUE)
            my_dq.append(val)
            std_dq.append(val)
            check_consistency("append")

        elif op_choice < 0.5:  # appendleft
            val = random.randint(0, MAX_VALUE)
            my_dq.appendleft(val)
            std_dq.appendleft(val)
            check_consistency("appendleft")

        elif op_choice < 0.6:  # pop
            if len(std_dq) > 0:
                res_my = my_dq.pop()
                res_std = std_dq.pop()
                assert res_my == res_std, "pop 返回值不一致"
                check_consistency("pop")
            else:  # 测试空队列 pop
                try:
                    my_dq.pop()
                    assert False, "空队列 pop 未抛出 IndexError"
                except IndexError:
                    pass

        elif op_choice < 0.7:  # popleft
            if len(std_dq) > 0:
                res_my = my_dq.popleft()
                res_std = std_dq.popleft()
                assert res_my == res_std, "popleft 返回值不一致"
                check_consistency("popleft")
            else:  # 测试空队列 popleft
                try:
                    my_dq.popleft()
                    assert False, "空队列 popleft 未抛出 IndexError"
                except IndexError:
                    pass

        elif op_choice < 0.8:  # __getitem__ 和 __setitem__
            if len(std_dq) > 0:
                idx = random.randint(0, len(std_dq) - 1)
                # 50% 概率测试负索引
                if random.random() < 0.5:
                    idx = idx - len(std_dq)

                # 测试 getitem
                assert my_dq[idx] == std_dq[idx], f"__getitem__ 在索引 {idx} 处不一致"

                # 测试 setitem
                if random.random() < 0.3:  # 降低 setitem 频率
                    val = random.randint(0, MAX_VALUE)
                    my_dq[idx] = val
                    std_dq[idx] = val
                    check_consistency(f"__setitem__ at index {idx}")

        elif op_choice < 0.85:  # reverse
            if len(std_dq) > 0:
                my_dq.reverse()
                std_dq.reverse()
                check_consistency("reverse")

        elif op_choice < 0.9:  # clear
            my_dq.clear()
            std_dq.clear()
            check_consistency("clear")

        else:  # 混合长序列操作，触发 rebalance
            for _ in range(random.randint(5, 20)):
                if random.random() < 0.5:
                    val = random.randint(0, MAX_VALUE)
                    my_dq.append(val)
                    std_dq.append(val)
                else:
                    if len(std_dq) > 0:
                        my_dq.popleft()
                        std_dq.popleft()
            check_consistency("混合长序列操作")

        if (i + 1) % (NUM_OPERATIONS // 10) == 0:
            print(f"  已完成 {(i + 1)} / {NUM_OPERATIONS} 次操作...")

    print("\n所有测试通过！MyDeque 与 collections.deque 行为一致。")


if __name__ == "__main__":
    run_tests()
