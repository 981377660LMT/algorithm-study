# 力扣改版：
# https://leetcode.cn/problems/count-complete-substrings/
# s中的每个字符恰好出现`k的倍数次`.
# 给你一个字符串 word 和一个整数 k 。
# 如果 word 的一个子字符串 s 满足以下条件，我们称它是 完全字符串：
# !s 中每个字符 恰好 出现 k 的倍数次。
# 相邻字符在字母表中的顺序 至多 相差 2 。也就是说，s 中两个相邻字符 c1 和 c2 ，它们在字母表中的位置相差 至多 为 2 。
# 请你返回 word 中 完全 子字符串的数目。
# 子字符串 指的是一个字符串中一段连续 非空 的字符序列。
# 每个数，出现mod k次时的哈希值相等

# !一种字符出现k次时异或为0

from collections import defaultdict
from random import randint
from fastHash import AllCountKChecker


class Solution:
    def countCompleteSubstrings(self, word: str, k: int) -> int:
        if k == 0:
            return 0

        H = AllCountKChecker(k)
        nums = [ord(x) - 97 for x in word]
        res = 0
        preXor = defaultdict(int, {0: 1})
        for i, num in enumerate(nums):
            if i > 0 and abs(nums[i] - nums[i - 1]) > 2:
                H.clear()
                preXor.clear()
                preXor[0] = 1
            H.add(num)
            hash_ = H.getHash()
            res += preXor[hash_]
            preXor[hash_] += 1
        return res


def bruteForce(word: str, k: int) -> int:
    res = 0
    for i in range(len(word)):
        for j in range(i, len(word)):
            curS = word[i : j + 1]
            counter = [0] * 26
            for char in word[i : j + 1]:
                counter[ord(char) - 97] += 1
            if all(x % k == 0 for x in counter) and all(
                abs(ord(x) - ord(y)) <= 2 for x, y in zip(curS, curS[1:])
            ):
                res += 1
    return res


if __name__ == "__main__":

    def check() -> None:
        for _ in range(100):
            word = "".join(chr(randint(97, 97 + 25)) for _ in range(randint(1, 10)))
            k = randint(1, 10)
            res1, res2 = bruteForce(word, k), Solution().countCompleteSubstrings(word, k)
            if res1 != res2:
                print(word, k, res1, res2)
                break
        else:
            print("ok")

    check()
    ...
