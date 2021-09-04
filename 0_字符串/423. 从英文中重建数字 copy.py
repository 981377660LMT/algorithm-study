import collections


class Solution:
    def originalDigits(self, s: str) -> str:
        c = collections.Counter(s)
        zero = c["z"]
        two = c["w"]
        four = c["u"]
        six = c["x"]
        eight = c["g"]
        one = c["o"] - zero - two - four
        three = c["t"] - eight - two
        five = c["f"] - four
        seven = c["s"] - six
        nine = c["i"] - six - eight - five
        res = ""
        for idx, val in enumerate([zero, one, two, three, four, five, six, seven, eight, nine]):
            print(idx, val)
            res += str(idx) * val
        return res


print(Solution().originalDigits("owoztneoer"))
