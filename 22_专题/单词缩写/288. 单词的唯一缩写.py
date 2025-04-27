# 单词的 缩写 需要遵循 <起始字母><中间字母数><结尾字母> 这样的格式。
# 如果单词只有两个字符，那么它就是它自身的 缩写 。

from typing import List


class ValidWordAbbr:
    __slots__ = "_abbr_map"

    def __init__(self, dictionary: List[str]):
        self._abbr_map = {}
        self._build(dictionary)

    def isUnique(self, word: str) -> bool:
        """
        A word is unique if:
          1. its abbreviation is not in the dictionary at all, or
          2. the only dictionary word with that abbreviation is itself.
        """
        abbr = self._make_abbr(word)
        if abbr not in self._abbr_map:
            return True
        return self._abbr_map[abbr] == word

    def _build(self, dictionary: List[str]) -> None:
        """
        Build the abbreviation map from the dictionary.
        """
        for word in dictionary:
            abbr = self._make_abbr(word)
            if abbr not in self._abbr_map:
                self._abbr_map[abbr] = word
            elif self._abbr_map[abbr] != word:
                self._abbr_map[abbr] = None

    @staticmethod
    def _make_abbr(word: str) -> str:
        """
        Produce the abbreviation:
          - if len <= 2: the word itself
          - else: first letter + (len-2) + last letter
        """
        n = len(word)
        if n <= 2:
            return word
        return f"{word[0]}{n-2}{word[-1]}"
