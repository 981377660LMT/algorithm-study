from typing import List
from collections import defaultdict


# 给你两个字符串数组 creators 和 ids ，和一个整数数组 views ，所有数组的长度都是 n 。
# 平台上第 i 个视频者是 creator[i] ，视频分配的 id 是 ids[i] ，且播放量为 views[i] 。

# 视频创作者的 流行度 是该创作者的 所有 视频的播放量的 总和 。
# 请找出流行度 最高 创作者以及该创作者播放量 最大 的视频的 id 。

# 如果存在多个创作者流行度都最高，则需要找出所有符合条件的创作者。
# 如果某个创作者存在多个播放量最高的视频，则只需要找出字典序最小的 id 。
# 返回一个二维字符串数组 answer ，其中 answer[i] = [creatori, idi]
# 表示 creatori 的流行度 最高 且其最流行的视频 id 是 idi ，可以按任何顺序返回该结果。


class Solution:
    def mostPopularCreator(
        self, creators: List[str], ids: List[str], views: List[int]
    ) -> List[List[str]]:
        userScore = defaultdict(int)
        userVideo = defaultdict(list)
        for user, videoId, view in zip(creators, ids, views):
            userScore[user] += view
            userVideo[user].append((videoId, view))
        maxScore = max(userScore.values())
        res = []
        for user in userScore:
            if userScore[user] == maxScore:
                bestVideo = min(userVideo[user], key=lambda x: (-x[1], x[0]))  # !按照播放量降序，字典序升序
                res.append([user, bestVideo[0]])
        return res
