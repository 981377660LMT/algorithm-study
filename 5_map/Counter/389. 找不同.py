import collections

# 字符串 t 由字符串 s 随机重排，然后在随机位置添加一个字母。
class Solution:
    def findTheDifference(self, s: str, t: str) -> str:
        return (collections.Counter(t) - collections.Counter(s)).popitem()[0]

