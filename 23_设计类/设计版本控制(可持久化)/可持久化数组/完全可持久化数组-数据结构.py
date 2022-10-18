# # 配列の初期状態
# A = [2, 0, 1, 9]
# # 更新操作
# operations = [('a.0', 2, 2, 'a.1'), ('a.1', 0, 1, 'a.2'), ('a.2', 1, 1, 'a.3'),
#               ('a.1', 3, 0, 'b.2'), ('a.3', 3, 3, 'a.4')]
# # クエリ
# queries = [('a.1', 2), ('a.4', 1), ('a.4', 3), ('b.2', 3)]
# https://qiita.com/wotsushi/items/72e7f8cdd674741ffd61#%E5%8F%82%E8%80%83%E8%A8%98%E4%BA%8B

# !完全可持久化数组


from typing import List
from PersistentArray2 import PersistentArray


class OnlineQuery:
    """持久化数组记录每个版本的数组"""

    def __init__(self, nums: List[int]) -> None:
        self.arr = PersistentArray.create(nums)
        self.git = dict({"a.0": self.arr})  # 保存每个版本的持久化数组

    def update(self, curVersion: str, index: int, value: int, nextVersion: str) -> None:
        """将curVersion版本的数组的index位置更新为value,得到nextVersion版本的数组"""
        preArr = self.git.get(curVersion, None)
        if preArr is None:
            raise Exception(f"版本{curVersion}不存在")
        curArr = preArr.update(index, value)
        self.git[nextVersion] = curArr

    def query(self, version: str, index: int) -> int:
        """查询version版本的数组的index位置的值"""
        arr = self.git.get(version, None)
        if arr is None:
            raise Exception(f"版本{version}不存在")
        return arr.get(index)


if __name__ == "__main__":
    persist = OnlineQuery([2, 0, 1, 9] * 2500)
    for _ in range(2000):
        persist.update("a.0", 2, 2, "a.1")
        persist.update("a.1", 0, 1, "a.2")
        persist.update("a.2", 1, 1, "a.3")
        persist.update("a.1", 3, 0, "b.2")
        persist.update("a.3", 3, 3, "a.4")
    for _ in range(2500):
        assert persist.query("a.1", 2) == 2
        assert persist.query("a.4", 1) == 1
        assert persist.query("a.4", 3) == 3
        assert persist.query("b.2", 3) == 0
