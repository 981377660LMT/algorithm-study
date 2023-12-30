# # 配列の初期状態
# A = [2, 0, 1, 9]
# # 更新操作
# operations = [('a.0', 2, 2, 'a.1'), ('a.1', 0, 1, 'a.2'), ('a.2', 1, 1, 'a.3'),
#               ('a.1', 3, 0, 'b.2'), ('a.3', 3, 3, 'a.4')]
# # クエリ
# queries = [('a.1', 2), ('a.4', 1), ('a.4', 3), ('b.2', 3)]
# https://qiita.com/wotsushi/items/72e7f8cdd674741ffd61#%E5%8F%82%E8%80%83%E8%A8%98%E4%BA%8B

# 倒序操作/反向操作 类似
# !5_map/经典题/prenext前驱后继/6092. 替换数组中的元素-prenext记录前驱后继.py
# 想象在一棵树上查询pre

from typing import List


class OnlineQuery:
    """在树上上跳查询最后一次更新的值"""

    def __init__(self, nums: List[int]) -> None:
        self.nums = nums
        self.pre = dict()

    def update(self, curVersion: str, index: int, value: int, nextVersion: str) -> None:
        """将curVersion版本的数组的index位置更新为value,得到nextVersion版本的数组"""
        self.pre[nextVersion] = (curVersion, index, value)

    def query(self, version: str, index: int) -> int:
        """查询version版本的数组的index位置的值

        倒着查询index下标最后一次被更新的值
        """
        curVersion = version
        while curVersion != "a.0":
            preVersion, pos, value = self.pre[curVersion]
            if pos == index:
                return value
            curVersion = preVersion
        return self.nums[index]


if __name__ == "__main__":
    preSearch = OnlineQuery([2, 0, 1, 9])
    preSearch.update("a.0", 2, 2, "a.1")
    preSearch.update("a.1", 0, 1, "a.2")
    preSearch.update("a.2", 1, 1, "a.3")
    preSearch.update("a.1", 3, 0, "b.2")
    preSearch.update("a.3", 3, 3, "a.4")
    assert preSearch.query("a.1", 2) == 2
    assert preSearch.query("a.4", 1) == 1
    assert preSearch.query("a.4", 3) == 3
    assert preSearch.query("b.2", 3) == 0
