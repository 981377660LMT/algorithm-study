# from collections import defaultdict
from heapq import heappop, heappush
from typing import Any, List


# 1. 因为要动态维护idPool最小值，所以要用堆/有序集合
# 2. 把 Video 封装成类可以增加代码可读性
# 3. 电影平台只要调用每部电影的api即可，不要把电影平台和电影紧耦合(命令模式)


class Video:
    def __init__(self, content: str):
        self.content = content
        self.platform = None
        self.id = -1
        self.view = 0
        self.like = 0
        self.dislike = 0

    def register(self, platform: 'VideoSharingPlatform') -> int:
        self.platform = platform
        self.id = heappop(self.platform.idPool)
        self.platform.idByVideo[self] = self.id
        self.platform.videoById[self.id] = self
        return self.id

    def destroy(self) -> int:
        if self.platform is None:
            return -1
        self.platform.idByVideo.pop(self)
        self.platform.videoById.pop(self.id)
        heappush(self.platform.idPool, self.id)
        self.platform = None
        return self.id

    def addView(self, count=1) -> int:
        self.view += count
        return self.view

    def addLike(self, count=1) -> int:
        self.like += count
        return self.like

    def addDislike(self, count=1) -> int:
        self.dislike += count
        return self.dislike

    def __getitem__(self, index: Any) -> str:
        return self.content[index]


class VideoSharingPlatform:
    """设计电影分享平台
    - 每部电影用数字字符串s表示,s[i]表示第i分钟的电影内容
    - 用户可以对电影点赞和点踩
    - 电影平台的管理人员需要统计每部电影的观看数、点赞数、点踩数
    - 电影的id是从0全局自增的,如果电影被删除,那么被删除的电影id可以被重新使用

    所有参数量级都是1e5
    """

    def __init__(self):
        self.idPool = list(range(int(1e5 + 10)))
        self.videoById = dict()
        self.idByVideo = dict()

    def upload(self, video: str) -> int:
        """上传电影,返回电影id"""
        newVideo = Video(video)
        return newVideo.register(self)

    def remove(self, videoId: int) -> None:
        """删除电影"""
        if videoId in self.videoById:
            self.videoById[videoId].destroy()

    def watch(self, videoId: int, startMinute: int, endMinute: int) -> str:
        """看电影,如果这部电影存在,那么观看数+1并返回看的这一段电影内容(闭区间),否则返回"-1"""
        if videoId not in self.videoById:
            return '-1'
        self.videoById[videoId].addView(1)
        return self.videoById[videoId][startMinute : endMinute + 1]

    def like(self, videoId: int) -> None:
        """点赞,如果这部电影存在"""
        if videoId in self.videoById:
            self.videoById[videoId].addLike(1)

    def dislike(self, videoId: int) -> None:
        """点踩,如果这部电影存在"""
        if videoId in self.videoById:
            self.videoById[videoId].addDislike(1)

    def getLikesAndDislikes(self, videoId: int) -> List[int]:
        """返回这部电影的[点赞数,点踩数];如果这部电影不存在,返回[-1]"""
        if videoId not in self.videoById:
            return [-1]
        return [self.videoById[videoId].like, self.videoById[videoId].dislike]

    def getViews(self, videoId: int) -> int:
        """返回这部电影的观看数;如果这部电影不存在,返回-1"""
        if videoId not in self.videoById:
            return -1
        return self.videoById[videoId].view


if __name__ == '__main__':
    platForm = VideoSharingPlatform()
    print(platForm.remove(0))
    print(platForm.watch(0, 0, 1))
    print(platForm.like(0))
    print(platForm.dislike(0))
    print(platForm.getLikesAndDislikes(0))
    print(platForm.getViews(0))

    # class Foo:
    #     def __init__(self):
    #         self.a = 'asas'

    #     def __getitem__(self, item):
    #         print(type(item))
    #         return self.a[item]

    # print(Foo()[1:3])
