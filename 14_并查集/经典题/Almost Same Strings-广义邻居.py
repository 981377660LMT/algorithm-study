# 字符串数组中 是否存在两个字符串仅仅相差一个字符
class Solution:
    def solve(self, words):
        visited = set()
        for word in words:
            for i in range(len(word)):
                if word[:i] + '*' + word[i + 1 :] in visited:
                    return True
                visited.add(word[:i] + '*' + word[i + 1 :])
        return False


# word[:i] + '*' + word[i + 1 :] 就是原字符串的广义邻居
