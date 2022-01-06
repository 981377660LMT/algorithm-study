from collections import defaultdict
from typing import List


class MagicDictionary(object):
    def __init__(self):
        self.adjMap = defaultdict(set)

    def buildDict(self, dictionary: List[str]) -> None:
        for word in set(dictionary):
            for i in range(len(word)):
                mode = word[:i] + '*' + word[i + 1 :]
                self.adjMap[mode].add(word)

    def search(self, searchWord: str) -> bool:
        for i in range(len(searchWord)):
            mode = searchWord[:i] + '*' + searchWord[i + 1 :]
            if len(self.adjMap[mode]) > 1:
                return True
            if len(self.adjMap[mode]) == 1 and searchWord not in self.adjMap[mode]:
                return True
        return False


if __name__ == '__main__':
    m = MagicDictionary()
    m.buildDict(["hello", "leetcode"])
    print(m.__dict__)

