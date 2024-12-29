# https://leetcode.cn/contest/weekly-contest-430/problems/find-the-lexicographically-largest-string-from-the-box-i/
# 100508. 从盒子中找出字典序最大的字符串
#
# !字典序最大后缀

from 最小表示法 import findSuffix


class Solution:
    def answerString(self, word: str, numFriends: int) -> str:
        if numFriends == 1:
            return word
        maxLex = findSuffix(word, isMin=False)
        len_ = len(word) - numFriends + 1
        return maxLex[:len_]
