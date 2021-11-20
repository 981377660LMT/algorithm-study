from typing import List
from operator import itemgetter

d = {
    'a': '.-',
    'b': '-...',
    'c': '-.-.',
    'd': '-..',
    'e': '.',
    'f': '..-.',
    'g': '--.',
    'h': '....',
    'i': '..',
    'j': '.---',
    'k': '-.-',
    'l': '.-..',
    'm': '--',
    'n': '-.',
    'o': '---',
    'p': '.--.',
    'q': '--.-',
    'r': '.-.',
    's': '...',
    't': '-',
    'u': '..-',
    'v': '...-',
    'w': '.--',
    'x': '-..-',
    'y': '-.--',
    'z': '--..',
}


class Solution:
    def uniqueMorseRepresentations(self, words: List[str]) -> int:
        return len(set("".join(itemgetter(*word)(d)) for word in words))


# print(itemgetter('abc')(d))
print(itemgetter(*'abc')(d))
