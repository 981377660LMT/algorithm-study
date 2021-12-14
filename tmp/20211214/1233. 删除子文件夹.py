from typing import List

# 你的任务是要删除该列表中的所有 子文件夹，并以 任意顺序 返回剩下的文件夹。

# 排序即可，第一个不被删
class Solution:
    def removeSubfolders(self, folder: List[str]) -> List[str]:
        res = []
        prefix = '$'
        for file in sorted(folder):
            if not file.startswith(prefix):
                res.append(file)
                prefix = file + '/'
        return res


print(Solution().removeSubfolders(["/a", "/a/b", "/c/d", "/c/d/e", "/c/f"]))
# 输出：["/a","/c/d","/c/f"]
# 解释："/a/b/" 是 "/a" 的子文件夹，而 "/c/d/e" 是 "/c/d" 的子文件夹
