# s 是一个 有效 的电子邮件或者电话号码
# 要想隐藏电子邮件地址中的个人信息：

# 名字 和 域名 部分的大写英文字母应当转换成小写英文字母。
# 名字 中间的字母（即，除第一个和最后一个字母外）必须用 5 个 "*****" 替换。

# 处理邮箱：找到@，然后利用index提取关键信息
# 处理手机号:提取数字，再用长度-10判断国家代码位数，查表，再加上公共部分
class Solution:
    def maskPII(self, S: str) -> str:
        at = S.find('@')
        if at >= 0:
            return (S[0] + "*" * 5 + S[at - 1 :]).lower()
        S = "".join(i for i in S if i.isdigit())
        return ["", "+*-", "+**-", "+***-"][len(S) - 10] + "***-***-" + S[-4:]


print(Solution().maskPII(S="LeetCode@LeetCode.com"))
# 输出："l*****e@leetcode.com"
# 解释：s 是一个电子邮件地址。
# 名字和域名都转换为小写，名字的中间用 5 个 * 替换。
