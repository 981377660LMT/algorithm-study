class Solution:
    def reverseOnlyLetters(self, s: str) -> str:
        sb = list(s)
        left, right = 0, len(s) - 1
        while left < right:
            while left < right and not s[left].isalpha():
                left += 1
            while left < right and not s[right].isalpha():
                right -= 1
            if left < right:
                sb[left], sb[right] = s[right], s[left]
                left += 1
                right -= 1

        return "".join(sb)


print(Solution().reverseOnlyLetters("a-bC-dEf-ghIj"))
