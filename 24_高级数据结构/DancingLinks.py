# https://atcoder.jp/contests/arc039/tasks/arc039_c
# https://mugen1337.github.io/procon/DataStructure/DancingLinks.hpp


from typing import Literal, Tuple


class DancingLinks:
    __slots__ = ("_x", "_y", "_mp")

    def __init__(self, sx=0, sy=0):
        self._x = sx
        self._y = sy
        self._mp = dict()
        self._update()

    def move(self, dir: Literal["L", "R", "U", "D", "l", "r", "u", "d"]):
        """沿着方向'L'/ 'R'/ 'U'/ 'D'移动，直到遇到未访问的位置为止
        L/l:x--
        R/r:x++
        U/u:y++
        D/d:y--
        """
        if dir == "R" or dir == "r":
            self._x = self._mp[(self._x, self._y, "R")]
        if dir == "L" or dir == "l":
            self._x = self._mp[(self._x, self._y, "L")]
        if dir == "U" or dir == "u":
            self._y = self._mp[(self._x, self._y, "U")]
        if dir == "D" or dir == "d":
            self._y = self._mp[(self._x, self._y, "D")]
        self._update()

    def get(self) -> Tuple[int, int]:
        return (self._x, self._y)

    def _update(self):
        state1 = (self._x, self._y, "R")
        state2 = (self._x, self._y, "L")
        state3 = (self._x, self._y, "U")
        state4 = (self._x, self._y, "D")
        if state1 not in self._mp:
            self._mp[state1] = self._x + 1
        if state2 not in self._mp:
            self._mp[state2] = self._x - 1
        if state3 not in self._mp:
            self._mp[state3] = self._y + 1
        if state4 not in self._mp:
            self._mp[state4] = self._y - 1

        self._mp[(self._mp[state1], self._y, "L")] = self._mp[state2]
        self._mp[(self._mp[state2], self._y, "R")] = self._mp[state1]
        self._mp[(self._x, self._mp[state3], "D")] = self._mp[state4]
        self._mp[(self._x, self._mp[state4], "U")] = self._mp[state3]


if __name__ == "__main__":
    # dl = DancingLinks(0, 0)
    # dl.move("R")
    # dl.move("R")
    # print(dl.get())
    # dl.move("U")
    # print(dl.get())
    # dl.move("L")
    # print(dl.get())
    # dl.move("D")
    # print(dl.get())

    # https://atcoder.jp/contests/arc039/tasks/arc039_c

    k = int(input())
    dirs = input()
    dl = DancingLinks()
    for d in dirs:
        dl.move(d)  # type: ignore
    print(*dl.get())
