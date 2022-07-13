from typing import List, Mapping, Tuple

# 1 <= s.length <= 10^5
# 给你一个只包含小写字母的字符串s ，你需要找到s中最多数目的非空子字符串，满足如下条件:
# 1.这些字符串之间互不重叠，也就是说对于任意两个子字符串s[i..j]和s[k..1]，要么j<k要么i> l。
# ⒉.如果一个子字符串包含字符char，那么s中所有char字符都应该在这个子字符串中。
# 请你找到满足上述条件的最多子字符串数目。如果有多个解法有相同的子字符串数目，请返回这些了字符串总长度最小的一个解。可以证明最小总长度解是唯一的。
# 请注意，你可以以任意顺序返回最优解的子字符串。


Interval = Tuple[int, int]


class Solution:
    def maxNumOfSubstrings(self, s: str) -> List[str]:
        """
        https://leetcode.com/problems/maximum-number-of-non-overlapping-substrings/discuss/743402/21-lines-Python-greedy-solution
        1.处理相同字母开头结尾的区间(最多26个)
        2.用尽可能多的线段去覆盖这个数轴，且线段间互不相交，线段之和最小
        选区间问题：贪心，尽量选结束时间早的区间，这样就剩下更多位置给右边
        处理出区间后判断每个区间是否合法
        从前往后遍历线段，每次遇到可以加入答案的线段，就贪心地将其加入答案数组即可
        """
        # 1.预处理区间
        intervalMap: Mapping[str, Interval] = dict()
        for char in set(s):
            start, end = s.find(char), s.rfind(char)
            intervalMap[char] = (start, end)

        # 2.对`每个字符`对应的区间寻找符合题意的边界
        validIntervals: List[Interval] = []
        for (start, end) in intervalMap.values():
            startCand, endCand = start, end
            while True:
                charInInterval = set(s[startCand : endCand + 1])
                for char in charInInterval:
                    startCand = min(startCand, intervalMap[char][0])
                    endCand = max(endCand, intervalMap[char][1])
                if (startCand, endCand) == (start, end):
                    break
                start, end = startCand, endCand
            validIntervals.append((start, end))

        # 3.排序，贪心选择结束早且更短的区间，类似leetcode452、646、1353
        validIntervals.sort(key=lambda x: (x[1], x[1] - x[0]))
        res, preEnd = [], -1
        for start, end in validIntervals:
            if start >= preEnd:
                res.append(s[start : end + 1])
                preEnd = end
        return res


# print(Solution().maxNumOfSubstrings(s="adefaddaccc"))
# print(Solution().maxNumOfSubstrings(s="abbaccd"))
# print(Solution().maxNumOfSubstrings(s="abab"))
print(Solution().maxNumOfSubstrings(s="abaabbcaaabbbccd"))
# ["d","abaabbcaaabbbcc"]
