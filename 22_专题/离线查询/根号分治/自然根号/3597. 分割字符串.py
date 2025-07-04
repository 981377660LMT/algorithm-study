# 3597. 分割字符串
# https://leetcode.cn/problems/partition-string/
# 给你一个字符串 s，按照以下步骤将其分割为 互不相同的段 ：
#
# 从下标 0 开始构建一个段。
# 逐字符扩展当前段，直到该段之前未曾出现过。
# 只要当前段是唯一的，就将其加入段列表，标记为已经出现过，并从下一个下标开始构建新的段。
# 重复上述步骤，直到处理完整个字符串 s。
# 返回字符串数组 segments，其中 segments[i] 表示创建的第 i 段。


from typing import List


class Solution:
    def partitionString2(self, s: str) -> List[str]:
        # !O(nsqrt(n))
        res = []
        visited = set()
        cur = ""
        for c in s:
            cur += c
            if cur not in visited:
                visited.add(cur)
                res.append(cur)
                cur = ""
        return res

    def partitionString(self, s: str) -> List[str]:
        trie = dict()
        cur = trie
        left = 0
        res = []
        for i, c in enumerate(s):
            if c in cur:
                cur = cur[c]
            else:
                cur[c] = dict()
                res.append(s[left : i + 1])
                left = i + 1
                cur = trie
        return res
