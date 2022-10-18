# # 配列の初期状態
# A = [2, 0, 1, 9]
# # 更新操作
# operations = [('a.0', 2, 2, 'a.1'), ('a.1', 0, 1, 'a.2'), ('a.2', 1, 1, 'a.3'),
#               ('a.1', 3, 0, 'b.2'), ('a.3', 3, 3, 'a.4')]
# # クエリ
# queries = [('a.1', 2), ('a.4', 1), ('a.4', 3), ('b.2', 3)]
# https://qiita.com/wotsushi/items/72e7f8cdd674741ffd61#%E5%8F%82%E8%80%83%E8%A8%98%E4%BA%8B


from typing import List


class OnlineQuery:
    """全量保存每个版本的数组"""

    def __init__(self, nums: List[int]) -> None:
        self.nums = nums
        self.git = dict({"a.0": nums})  # 保存每个版本的数组

    def update(self, curVersion: str, index: int, value: int, nextVersion: str) -> None:
        """将curVersion版本的数组的index位置更新为value,得到nextVersion版本的数组"""
        preNums = self.git.get(curVersion, None)
        if preNums is None:
            raise Exception(f"版本{curVersion}不存在")
        newNums = [value if i == index else num for i, num in enumerate(preNums)]
        self.git[nextVersion] = newNums

    def query(self, version: str, index: int) -> int:
        """查询version版本的数组的index位置的值"""
        nums = self.git.get(version, None)
        if nums is None:
            raise Exception(f"版本{version}不存在")
        return nums[index]


if __name__ == "__main__":
    fullBackup = OnlineQuery([2, 0, 1, 9] * 2500)
    for _ in range(2000):
        fullBackup.update("a.0", 2, 2, "a.1")
        fullBackup.update("a.1", 0, 1, "a.2")
        fullBackup.update("a.2", 1, 1, "a.3")
        fullBackup.update("a.1", 3, 0, "b.2")
        fullBackup.update("a.3", 3, 3, "a.4")
    for _ in range(2500):
        assert fullBackup.query("a.1", 2) == 2
        assert fullBackup.query("a.4", 1) == 1
        assert fullBackup.query("a.4", 3) == 3
        assert fullBackup.query("b.2", 3) == 0
