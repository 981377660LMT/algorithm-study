import collections


# 两个n之间要相隔n可
# 剪枝：可以从最远的开始搜索
def friendSequence(n):
    def dfs(cand, pos, start):
        """
        :param cand: 剩余可以填充的个数
        :param pos: 最终排序结果
        :param start: 可填充的起始位置
        :return:
        """
        if 0 not in pos:
            res.append(pos[:])
            return
        for i in range(n, 0, -1):
            if cand[i] == 0:
                continue
            # 如果i剩余的值是大于0的
            while pos[start] != 0 and start < 2 * n:
                start += 1
            if (
                start < 2 * n
                and pos[start] == 0
                and start + i + 1 < 2 * n
                and pos[start + i + 1] == 0
            ):
                # 要求当前位置和与其间隔i个位置都没被填充，且注意索引没有越界
                pos[start], pos[start + i + 1] = i, i
                cand[i] -= 2
                dfs(cand, pos, start + 1)
                pos[start], pos[start + i + 1] = 0, 0
                cand[i] += 2

    pos = [0] * 2 * n
    res = []
    cand = collections.defaultdict(int)
    for i in range(1, n + 1):
        cand[i] = 2
    dfs(cand, pos, 0)
    # for sin_res in res:
    #     print(sin_res[::-1])
    return res


print(friendSequence(4))
print(friendSequence(15))
# n=15需要压缩状态+记忆化
