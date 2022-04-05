from collections import defaultdict
from typing import Counter, List


class Encrypter:
    def __init__(self, keys: List[str], values: List[str], dictionary: List[str]):
        self.keyToValue = defaultdict(str, {k: v for k, v in zip(keys, values)})
        self.counter = Counter()
        # 注意反解的时候 dictionary 里的字符 要在原keys中出现
        for word in dictionary:
            if not all(w in self.keyToValue for w in word):
                continue
            self.counter[self.encrypt(word)] += 1

    def encrypt(self, word1: str) -> str:
        return ''.join(self.keyToValue[c] for c in word1)

    def decrypt(self, word2: str) -> int:
        return self.counter[word2]


# Your Encrypter object will be instantiated and called as such:
# obj = Encrypter(keys, values, dictionary)
# param_1 = obj.encrypt(word1)
# param_2 = obj.decrypt(word2)
