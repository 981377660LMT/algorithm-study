"""因为rowId自增,删行不影响rowId,不用真的删行,只要把行的内容设为空就行"""

from typing import List


class SQL:
    def __init__(self, names: List[str], columns: List[int]):
        self._table = {name: [] for name in names}

    def insertRow(self, name: str, row: List[str]) -> None:
        self._table[name].append(row)

    def deleteRow(self, name: str, rowId: int) -> None:
        rowId -= 1
        self._table[name][rowId] = []

    def selectCell(self, name: str, rowId: int, columnId: int) -> str:
        rowId, columnId = rowId - 1, columnId - 1
        return self._table[name][rowId][columnId]


if __name__ == "__main__":
    # ["SQL", "insertRow", "selectCell", "insertRow", "deleteRow", "selectCell"]
    # [[["one", "two", "three"], [2, 3, 1]], ["two", ["first", "second", "third"]], ["two", 1, 3], ["two", ["fourth", "fifth", "sixth"]], ["two", 1], ["two", 2, 2]]
    sql = SQL(["one", "two", "three"], [2, 3, 1])
    sql.insertRow("two", ["first", "second", "third"])
    print(sql.selectCell("two", 1, 3))
    sql.insertRow("two", ["fourth", "fifth", "sixth"])
    sql.deleteRow("two", 1)
    print(sql.selectCell("two", 2, 2))
