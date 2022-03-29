class Solution:
    def solve(self, string):
        endswith = [1, 0, 0, 0]
        for char in string:
            index = ord(char) - ord('a') + 1
            # endswith[index]表示取前, endswith[index - 1]表示不取前
            endswith[index] += endswith[index] + endswith[index - 1]
        return endswith[-1]
