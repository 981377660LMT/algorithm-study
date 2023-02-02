from Trie import Trie as T


class Trie:
    def __init__(self):
        self.trie = T()

    def insert(self, word: str) -> None:
        self.trie.insert(word)

    def countWordsEqualTo(self, word: str) -> int:
        res = self.trie.countWord(word)
        return res[-1] if res else 0

    def countWordsStartingWith(self, prefix: str) -> int:
        res = self.trie.countWordStartsWith(prefix)
        return res[-1] if res else 0

    def erase(self, word: str) -> None:
        self.trie.remove(word)


# Your Trie object will be instantiated and called as such:
# obj = Trie()
# obj.insert(word)
# param_2 = obj.countWordsEqualTo(word)
# param_3 = obj.countWordsStartingWith(prefix)
# obj.erase(word)
