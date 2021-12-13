from sortedcontainers import SortedList
from itertools import count

# 一个观光景点由它的名字 name 和景点评分 score 组成，其中 name 是所有观光景点中 唯一 的字符串
# 景点评分 越高 ，这个景点越好。如果有两个景点的评分一样，那么 字典序较小 的景点更好。
# 你需要搭建一个系统，查询景点的排名


class SORTracker:
    def __init__(self):
        self.slist = SortedList()
        self.timer = count()

    # 添加 景点，每次添加 一个 景点。
    def add(self, name: str, score: int) -> None:
        self.slist.add((-score, name))

    # 查询 已经添加景点中第 i 好 的景点，其中 i 是系统目前位置查询的次数（包括当前这一次）。
    def get(self) -> str:
        return self.slist[next(self.timer)][1]


s = SORTracker()
s.add(*["bradford", 2])
s.add(*["alps", 2])
s.add(*["orlando", 3])
print(s.get())
print(s.get())
print(s.get())
print(s.get())
