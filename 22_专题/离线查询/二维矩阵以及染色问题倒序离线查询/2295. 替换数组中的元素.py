# 给定一个数组nums和一系列操作operations，每个操作包含两个元素[a,b]，
# 表示将数组中所有值为a的元素替换为b，请你返回最终结果数组。
# 注意原数组可能包含重复元素。(不含重复元素可以用哈希表记录index修改)
# !在线:更新前驱后继(整个版本是一个单链表)
# !离线:倒序+修改后继(整个版本是一个单链表)

# last0 -> last1 -> last2 -> last3

from typing import List


class Solution:
    def arrayChange(self, nums: List[int], operations: List[List[int]]) -> List[int]:
        """离线解法."""
        next_ = dict()
        for from_, to in operations[::-1]:  # from_ -> to 这个关系
            next_[from_] = next_.get(to, to)
        return [next_.get(num, num) for num in nums]

    def arrayChange2(self, nums: List[int], operations: List[List[int]]) -> List[int]:
        """在线解法."""
        pre, next_ = dict(), dict()  # !模拟双向链表,更新时修改前驱后继
        for from_, to in operations:
            pre[to] = pre.get(from_, from_)
            next_[pre[to]] = to
        return [next_.get(num, num) for num in nums]
