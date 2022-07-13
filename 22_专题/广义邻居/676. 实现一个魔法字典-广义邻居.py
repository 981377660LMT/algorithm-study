from collections import defaultdict
from typing import List

# 将searchWord替换一个字母后，是否存在于dict
# 广义邻居


class MagicDictionary(object):
    def __init__(self):
        self.adjMap = defaultdict(set)  # 广义邻居

    def buildDict(self, dictionary: List[str]) -> None:
        for word in set(dictionary):
            for i in range(len(word)):
                replaced = word[:i] + "*" + word[i + 1 :]
                self.adjMap[replaced].add(word)

    def search(self, searchWord: str) -> bool:
        """给定一个字符串 searchWord ，判定能否只将字符串中一个字母换成另一个字母，使得所形成的新字符串能够与字典中的任一字符串匹配"""
        for i in range(len(searchWord)):
            replaced = searchWord[:i] + "*" + searchWord[i + 1 :]
            if len(self.adjMap[replaced]) > 1:
                return True

            # 唯一的邻居不是自己
            if len(self.adjMap[replaced]) == 1 and searchWord not in self.adjMap[replaced]:
                return True

        return False


if __name__ == "__main__":
    m = MagicDictionary()
    m.buildDict(["hello", "leetcode"])
    print(m.__dict__)
