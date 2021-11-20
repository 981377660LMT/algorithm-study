# for else 的用法
from typing import List


def countCharacters(self, words: List[str], chars: str) -> int:
    res = 0
    for word in words:
        for char in word:
            if word.count(char) > chars.count(char):
                break
        else:
            res += len(word)
    return res


# 假如你可以用 chars 中的『字母』（字符）拼写出 words 中的某个『单词』（字符串），那么我们就认为你掌握了这个单词。
