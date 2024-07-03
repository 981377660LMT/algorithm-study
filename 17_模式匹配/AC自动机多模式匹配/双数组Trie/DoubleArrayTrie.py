class DoubleArrayTrie:
    __slots__ = ("_base", "_check", "_size")

    def __init__(self):
        self._base = [0]
        self._check = [-1]
        self._size = 1

    def insert(self, word: str) -> None:
        pos = 0
        for c in word:
            index = ord(c) + 1
            if pos + index >= self._size or self._check[pos + index] != pos:
                self._expand(pos, index)
            pos = self._base[pos] + index
        self._base[pos] = -len(word)

    def search(self, word: str) -> bool:
        pos = 0
        for c in word:
            index = ord(c) + 1
            if pos + index >= self._size or self._check[pos + index] != pos:
                return False
            pos = self._base[pos] + index
        return self._base[pos] == -len(word)

    def delete(self, word: str) -> bool:
        if not self.search(word):
            return False
        pos = 0
        for c in word:
            idx = ord(c) + 1
            if self._base[pos] + idx == -len(word):
                self._base[pos] = 0
            pos = self._base[pos] + idx
        return True

    def _expand(self, cur_pos: int, idx: int) -> None:
        max_idx = max(idx, self._size - cur_pos)
        self._base.extend([0] * max_idx)
        self._check.extend([-1] * max_idx)
        self._size += max_idx
        self._check[cur_pos + idx] = cur_pos


if __name__ == "__main__":
    trie = DoubleArrayTrie()
    trie.insert("apple")
    print(trie.search("a"))
