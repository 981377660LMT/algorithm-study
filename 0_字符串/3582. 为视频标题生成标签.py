# 3582. 为视频标题生成标签
# https://leetcode.cn/problems/generate-tag-for-video-caption/description/
# 给你一个字符串 caption，表示一个视频的标题。
# 需要按照以下步骤 按顺序 生成一个视频的 有效标签 ：
# 将 所有单词 组合为单个 驼峰命名字符串 ，并在前面加上 '#'。驼峰命名字符串 指的是除第一个单词外，其余单词的首字母大写，且每个单词的首字母之后的字符必须是小写。
# 移除 所有不是英文字母的字符，但 保留 第一个字符 '#'。
# 将结果 截断 为最多 100 个字符。
# 对 caption 执行上述操作后，返回生成的 标签 。


class Solution:
    def generateTag(self, caption: str) -> str:
        res = ["#"]
        words = caption.split()
        for i, w in enumerate(words):
            if i == 0:
                res.append(w.lower())
            else:
                res.append(w.capitalize())
        return "".join(res)[:100]
