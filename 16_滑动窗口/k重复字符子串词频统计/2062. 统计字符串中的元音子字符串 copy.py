class Solution:
    def countVowelSubstrings(self, word: str) -> int:
        res, n = 0, len(word)
        vowels = set("aeiou")
        for left in range(n - 4):
            for right in range(left + 5, n + 1):
                if set(word[left:right]) == vowels:
                    res += 1
        return res
