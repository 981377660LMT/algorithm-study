from typing import List
from collections import defaultdict, deque

# watchedVideos[i]  和 friends[i] 分别表示 id = i 的人观看过的视频列表和他的好友列表。
# Level 为 k 的视频包含所有从你出发，最短距离为 k 的好友观看过的视频。
# 给定你的 id  和一个 level 值，请你找出所有指定 level 的视频，并将它们按观看频率升序返回。如果有频率相同的视频，请将它们按字母顺序从小到大排列。
class Solution:
    def watchedVideosByFriends(
        self, watchedVideos: List[List[str]], friends: List[List[int]], id: int, level: int
    ) -> List[str]:
        # ------------------------- 记忆化+bfs 波纹法
        queue = deque([id])
        visited = set([id])

        L = 0
        while L < level:
            cur_len = len(queue)
            for _ in range(cur_len):
                xID = queue.popleft()
                for yID in friends[xID]:
                    if yID not in visited:
                        visited.add(yID)
                        queue.append(yID)
            L += 1

        # ---- 统计
        video_freq = defaultdict(int)
        for ID in queue:
            for video in watchedVideos[ID]:
                video_freq[video] += 1
        # ---- 排序
        res = list(video_freq.keys())
        res.sort(key=lambda v: (video_freq[v], v))
        return res


print(
    Solution().watchedVideosByFriends(
        watchedVideos=[["A", "B"], ["C"], ["B", "C"], ["D"]],
        friends=[[1, 2], [0, 3], [0, 3], [1, 2]],
        id=0,
        level=1,
    )
)

