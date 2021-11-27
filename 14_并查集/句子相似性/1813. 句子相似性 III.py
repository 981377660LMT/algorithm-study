# 如果两个句子 sentence1 和 sentence2 ，
# 可以通过往其中一个句子插入一个任意的句子（可以是空句子）而得到另一个句子，
# 那么我们称这两个句子是 相似的

# 即判断:s1(短)是否由s2删除某个连续单词串得到：这个特点是短的单词比较头尾，可以全部去除
from collections import deque


class Solution:
    def areSentencesSimilar(self, s1: str, s2: str) -> bool:
        if len(s1) > len(s2):
            s1, s2 = s2, s1
        q1, q2 = deque(s1.split()), deque(s2.split())
        while q1 and q1[0] == q2[0]:
            q1.popleft()
            q2.popleft()
        while q1 and q1[-1] == q2[-1]:
            q1.pop()
            q2.pop()
        return not q1


print(Solution().areSentencesSimilar("My name is Haley", "My Haley"))

