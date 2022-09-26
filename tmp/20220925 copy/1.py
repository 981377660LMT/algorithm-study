from typing import List


class Solution:
    def sortPeople(self, names: List[str], heights: List[int]) -> List[str]:
        """按身高 降序 顺序返回对应的名字数组 names 。"""
        people = [(name, height) for name, height in zip(names, heights)]
        return [name for name, _ in sorted(people, key=lambda x: x[1], reverse=True)]
