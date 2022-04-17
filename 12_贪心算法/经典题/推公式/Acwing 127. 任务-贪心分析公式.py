# # 第 i 个任务的难度级别为 yi，完成任务所需时间为 xi 分钟。
# # 收入公式: 500 * xi + 2 * yi
# # 任务个数，机器数<=1e5
# # 0<xi<1440,
# # 0≤yi≤100

# # 1.贪心
# # x变动1, 收入变化500
# # y按照最大变动100, 收入变化才只有200.
# # 因此, 非常容易想到, 对于任务的优先级, 应该是先基于x, 再基于y
# # 2.逆序遍历任务.
# # 类似Acwing 110防晒的思路


from bisect import bisect_left, insort_left
from typing import Any, Generic, Iterable, Optional, Protocol, TypeVar, Union


def main() -> None:
    m, t = map(int, input().split())

    machines = []
    tasks = []

    for _ in range(m):
        timeLimit, level = map(int, input().split())
        machines.append((timeLimit, level))

    for _ in range(t):
        cost, level = map(int, input().split())
        tasks.append((cost, level))

    machines.sort()
    tasks.sort()
    maxCount, maxMoney = 0, 0

    cur = SortedList()
    mi = m - 1
    for ti in range(t - 1, -1, -1):
        while mi >= 0 and machines[mi][0] >= tasks[ti][0]:
            cur.add(machines[mi][1])  # 所有level都放入cur
            mi -= 1
        pos = cur.bisect_left(tasks[ti][1])  # 找到最差的机器
        if pos != len(cur):
            maxCount += 1
            maxMoney += 500 * tasks[ti][0] + 2 * tasks[ti][1]
            cur.pop(pos)

    print(maxCount, maxMoney)


# # 对于每个测试用例，输出两个整数，代表公司今天可以完成的最大任务数以及他们将获得的收入。
while True:
    try:
        main()
    except EOFError:
        break


class SupportsDunderLT(Protocol):
    def __lt__(self, __other: Any) -> bool:
        ...


class SupportsDunderGT(Protocol):
    def __gt__(self, __other: Any) -> bool:
        ...


S = TypeVar('S', bound=Union[SupportsDunderLT, SupportsDunderGT])


class SortedList(Generic[S]):
    def __init__(self, iterable: Optional[Iterable[S]] = None) -> None:
        self._list = []
        if iterable is not None:
            for item in iterable:
                self.add(item)

    def add(self, item: S) -> None:
        insort_left(self._list, item)

    def pop(self, index: int) -> S:
        return self._list.pop(index)

    def bisect_left(self, item: S) -> int:
        return bisect_left(self._list, item)

    def __getitem__(self, index: int) -> S:
        return self._list[index]

    def __len__(self) -> int:
        return len(self._list)


# 输入样例：
# 1 2
# 100 3
# 100 2
# 100 1

# 输出样例：
# 1 50004
