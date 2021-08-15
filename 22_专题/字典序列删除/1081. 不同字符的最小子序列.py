import collections

# 需要保证每个数最后正好出现一次
# 不多：使用visited保证最多出现一次
# 不少：使用counter防止多删


class Solution:
    def removeDuplicateLetters(self, s) -> str:
        stack = []
        visited = set()
        remain_counter = collections.Counter(s)
        for letter in s:
            if letter not in visited:
                while stack and stack[-1] > letter and remain_counter.get(stack[-1]):
                    visited.discard(stack.pop())
                stack.append(letter)
                visited.add(letter)
            remain_counter[letter] -= 1
        return ''.join(stack)


print(Solution().removeDuplicateLetters("cbacdcbc"))
