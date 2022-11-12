from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 message 和一个正整数 limit 。

# 你需要根据 limit 将 message 分割 成一个或多个 部分 。每个部分的结尾都是 "<a/b>" ，其中 "b" 用分割出来的总数 替换， "a" 用当前部分所在的编号 替换 ，编号从 1 到 b 依次编号。除此以外，除了最后一部分长度 小于等于 limit 以外，其他每一部分（包括结尾部分）的长度都应该 等于 limit 。

# 你需要确保分割后的结果数组，删掉每部分的结尾并 按顺序 连起来后，能够得到 message 。同时，结果数组越短越好。

# 请你返回 message  分割后得到的结果数组。如果无法按要求分割 message ，返回一个空数组。
class Solution:
    def splitMessage(self, message: str, limit: int) -> List[str]:
        def check(group: int) -> List[int]:
            """能否分成group组"""
            L = len(str(group)) + 3  # !还要加上当前的分组
            wordCount = 0
            for i in range(group - 1):
                sufLen = len(str(i + 1)) + L
                if sufLen > limit:
                    return []
                wordCount += limit - sufLen
            remainCount = n - wordCount
            if remainCount + L + len(str(group)) > limit:
                return []

            # 合理的
            res = []
            preI = 0
            for i in range(group - 1):
                sufLen = len(str(i + 1)) + L
                wordLen = limit - sufLen
                res.append(f"{message[preI : preI + wordLen]}<{i + 1}/{group}>")
                preI += wordLen

            # 最后一组
            res.append(f"{message[preI:]}<{group}/{group}>")
            return res

        n = len(message)
        # 直接枚举分几组
        for i in range(1, n + 1):
            res = check(i)
            if res:
                return res
        return []


print(Solution().splitMessage(message="this is really a very awesome message", limit=9))
