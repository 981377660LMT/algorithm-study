from collections import defaultdict
from typing import List

# 将searchWord替换一个字母后，是否存在于dict
# 广义邻居


class MagicDictionary(object):
    def __init__(self):
        self.adjMap = defaultdict(set)

    def buildDict(self, dictionary: List[str]) -> None:
        for word in set(dictionary):
            for i in range(len(word)):
                replace = word[:i] + '*' + word[i + 1 :]
                self.adjMap[replace].add(word)

    def search(self, searchWord: str) -> bool:
        for i in range(len(searchWord)):
            replace = searchWord[:i] + '*' + searchWord[i + 1 :]
            if len(self.adjMap[replace]) > 1:
                return True
            if len(self.adjMap[replace]) == 1 and searchWord not in self.adjMap[replace]:
                return True
        return False


if __name__ == '__main__':
    m = MagicDictionary()
    m.buildDict(["hello", "leetcode"])
    print(m.__dict__)

