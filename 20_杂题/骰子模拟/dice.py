# https://maspypy.github.io/library/other/dice.hpp
# 模拟骰子
# 次の番号付けるに従う：UFRLBD
# i, 5-i が反対の面になっている (i=0,1,2,3,4,5)


from typing import Generic, List, Literal, Tuple, TypeVar


T = TypeVar("T", int, str)
Faces = Tuple[T, T, T, T, T, T]


class Dice(Generic[T]):
    """
    https://maspypy.github.io/library/other/dice.hpp を参考にした.

    - 次の番号付けるに従う:UFRLBD.
    - i, 5-i が反対の面になっている (i=0,1,2,3,4,5).
    """

    __slots__ = ("_f", "_all_faces", "_hash")

    def __init__(self, faces: "Faces[T]"):
        self._f = faces
        self._all_faces = None
        self._hash = None

    def rotate_up_to(self, direction: Literal["F", "R", "L", "B"]) -> None:
        """U のうつる先となる FRLB を指定する."""
        if direction == "F":
            self._f = (self._f[4], self._f[0], self._f[2], self._f[3], self._f[5], self._f[1])
        elif direction == "R":
            self._f = (self._f[3], self._f[1], self._f[0], self._f[5], self._f[4], self._f[2])
        elif direction == "L":
            self._f = (self._f[2], self._f[1], self._f[5], self._f[0], self._f[4], self._f[3])
        elif direction == "B":
            self._f = (self._f[1], self._f[5], self._f[2], self._f[3], self._f[0], self._f[4])

    def get_all_faces(self) -> List["Faces[T]"]:
        """24種の面の組み合わせを返す."""
        if self._all_faces is not None:
            return self._all_faces  # type: ignore
        res = [None] * 24
        tmp = [None] * 24
        tmp[:4] = [(0, 1, 2), (0, 4, 3), (5, 1, 3), (5, 4, 2)]  # type: ignore
        for i in range(4):
            a, b, c = tmp[i]  # type: ignore
            tmp[4 + i] = (b, c, a)  # type: ignore
            tmp[8 + i] = (c, a, b)  # type: ignore
        for i in range(12):
            a, b, c = tmp[i]  # type: ignore
            tmp[12 + i] = (5 - b, a, c)  # type: ignore
        for i in range(24):
            a, b, c = tmp[i]  # type: ignore
            res[i] = (self._f[a], self._f[b], self._f[c], self._f[~c], self._f[~b], self._f[~a])  # type: ignore
        self._all_faces = res
        return res  # type: ignore

    @property
    def up(self) -> T:
        return self._f[0]

    @property
    def front(self) -> T:
        return self._f[1]

    @property
    def right(self) -> T:
        return self._f[2]

    @property
    def left(self) -> T:
        return self._f[3]

    @property
    def back(self) -> T:
        return self._f[4]

    @property
    def down(self) -> T:
        return self._f[5]

    def __repr__(self) -> str:
        return f"Dice(U: {self.up}, F: {self.front}, R: {self.right}, L: {self.left}, B: {self.back}, D: {self.down})"

    def __getitem__(self, i: Literal[0, 1, 2, 3, 4, 5]) -> T:
        return self._f[i]

    def __hash__(self) -> int:
        if self._hash is not None:
            return self._hash
        min_f = min(self.get_all_faces())
        self._hash = hash(min_f)
        return self._hash

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, Dice):
            return False
        return self.__hash__() == other.__hash__()


if __name__ == "__main__":
    # https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=ITP1_11_D&lang=ja
    # 给定n个骰子,判断他们是不是都不一样
    n = int(input())
    S = set()
    for _ in range(n):
        faces = tuple(map(int, input().split()))
        d = Dice(faces)
        S.add(d)
    print("Yes" if len(S) == n else "No")
