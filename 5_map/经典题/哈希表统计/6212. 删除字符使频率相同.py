"""是否可以 删除一个字母后剩余字母频率相等"""


from collections import Counter


# !如果删除一个字母后，word 中剩余所有字母的出现频率都相同，
# !那么返回 true ，否则返回 false 。
class Solution:
    def equalFrequency(self, word: str) -> bool:
        """暴力法"""
        for i in range(len(word)):
            s = word[:i] + word[i + 1 :]
            if len(set(Counter(s).values())) == 1:
                return True
        return False

    def equalFrequency2(self, word: str) -> bool:
        counter = Counter(word)
        freq = sorted(counter.values())

        # !一种字母
        if len(freq) == 1:
            return True
        # !删最少的 1 3 3 3
        if freq[0] == 1 and freq[1] == freq[-1]:
            return True
        # !删最多的 3 3 3 4
        if freq[0] == freq[-2] and freq[-1] == freq[-2] + 1:
            return True
        return False


print(Solution().equalFrequency("aazz"))
print(Solution().equalFrequency2("aazz"))
