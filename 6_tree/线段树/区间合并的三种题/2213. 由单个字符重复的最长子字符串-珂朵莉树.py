# https://leetcode-cn.com/problems/longest-substring-of-one-repeating-character/solution/python-guo-ran-wo-huan-shi-geng-xi-huan-olhop/

from itertools import groupby
from typing import List
from sortedcontainers import SortedList, SortedDict

# 1 <= s.length <= 1e5
# 1 <= k <= 1e5


# 1.果然，我还是喜欢sortedList啊
# 2.复杂的操作要提取成函数，思路清晰一些
# 3. 求当前所在区间用pos=sMap.bisect_right(index) - 1 ，前一个区间pos-1 后一个区间pos+1


# 珂朵莉树的适用范围是有区间赋值操作且数据随机的题目。其实珂朵莉树看上去并不像是树状数据结构，但因为一般要用到std::set，而std::set是用红黑树实现的，
# 所以也不算名不副实。在随机数据下，珂朵莉树可以达到 nloglogn 的复杂度


class Solution:
    def longestRepeating(self, s: str, queryCharacters: str, queryIndices: List[int]) -> List[int]:
        """
        第 i 个查询会将 s 中位于下标 queryIndices[i] 的字符更新为 queryCharacters[i] 。
        返回一个长度为 k 的数组 lengths ，其中 lengths[i] 是在执行第 i 个查询 之后 s 中仅由 单个字符重复 组成的 最长子字符串 的 长度 。
        """

        def split(index: int) -> None:
            """
            向左断开区间；如果 index 不是某一个区间的起点，则将 index 所在区间分成两段，第二段以 index 为起点
            断开index与左右两边的连接用 `split(index),split(index+1)`
            如果是对[left,right]操作，则需要先split(left),split(right+1)
            """
            if index < 0 or index >= n:
                return

            curPos = sMap.bisect_right(index) - 1
            s1, e1 = sMap.peekitem(curPos)
            if s1 == index:
                return
            sMap.popitem(curPos)
            sMap[s1] = index - 1
            sMap[index] = e1
            sList.remove(e1 - s1 + 1)
            sList.add(index - 1 - s1 + 1)
            sList.add(e1 - index + 1)

        def union(index: int) -> None:
            """
            向左合并区间；如果以 index 为起点的区间和其前一个区间内的字符相同，合并两个区间
            连接index与左右两边的区间用 `union(index),union(index+1)`
            如果是对[left,right]操作，则需要先union(left),union(right+1)
            """
            if index < 0 or index >= n:
                return

            curPos = sMap.bisect_right(index) - 1
            prePos = curPos - 1
            if prePos < 0:
                return
            (s1, e1), (s2, e2) = sMap.peekitem(prePos), sMap.peekitem(curPos)
            if chars[s1] == chars[s2]:
                sMap.popitem(curPos)
                sMap[s1] = e2
                sList.remove(e2 - s2 + 1)
                sList.remove(e1 - s1 + 1)
                sList.add(e2 - s1 + 1)

        sMap = SortedDict()  # start=>end
        sList = SortedList(key=lambda x: -x)

        # 1.初始化
        start = 0
        for _, group in groupby(s):
            len_ = len(list(group))
            sMap[start] = start + len_ - 1
            sList.add(len_)
            start += len_

        res = [1] * len(queryIndices)
        n, chars = len(s), list(s)
        for i, (qc, qi) in enumerate(zip(queryCharacters, queryIndices)):
            if chars[qi] == qc:
                res[i] = sList[0]
                continue

            chars[qi] = qc

            # 断开qi
            split(qi)
            split(qi + 1)

            # 向左连接
            union(qi)
            union(qi + 1)

            res[i] = sList[0]

        return res


print(Solution().longestRepeating(s="babacc", queryCharacters="bcb", queryIndices=[1, 3, 3]))
# [3,3,4]

