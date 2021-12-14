from typing import List
from random import choice


class Master:
    def guess(self, word: str) -> int:
        ...


# 我们给出了一个由一些独特的单词组成的单词列表，每个单词都是 6 个字母长，并且这个列表中的一个单词将被选作秘密。
# 你可以调用 master.guess(word) 来猜单词。你所猜的单词应当是存在于原列表并且由 6 个小写字母组成的类型字符串。
# 此函数将会返回一个整型数字，表示你的猜测与秘密单词的准确匹配（值和位置同时匹配）的数目。此外，如果你的猜测不在给定的单词列表中，它将返回 -1。
class Solution:
    def findSecretWord(self, wordlist: List[str], master: 'Master') -> None:
        for _ in range(10):
            guess = choice(wordlist)
            hit = master.guess(guess)
            wordlist = [w for w in wordlist if sum(i == j for i, j in zip(guess, w)) == hit]


# 输入: secret = "acckzz", wordlist = ["acckzz","ccbazz","eiowzz","abcczz"]

# 解释:

# master.guess("aaaaaa") 返回 -1, 因为 "aaaaaa" 不在 wordlist 中.
# master.guess("acckzz") 返回 6, 因为 "acckzz" 就是秘密，6个字母完全匹配。
# master.guess("ccbazz") 返回 3, 因为 "ccbazz" 有 3 个匹配项。
# master.guess("eiowzz") 返回 2, 因为 "eiowzz" 有 2 个匹配项。
# master.guess("abcczz") 返回 4, 因为 "abcczz" 有 4 个匹配项。

# 我们调用了 5 次master.guess，其中一次猜到了秘密，所以我们通过了这个测试用例。

