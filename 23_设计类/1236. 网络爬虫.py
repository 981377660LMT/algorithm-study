from typing import List, Protocol
from collections import deque


class HtmlParser(Protocol):
    def getUrls(self, url: str) -> List[str]:
        raise NotImplementedError('Must implement getUrls()')


# 请你实现一个网络爬虫，以实现爬取同 startUrl 拥有相同 域名标签 的全部链接。该爬虫得到的全部链接可以 任何顺序 返回结果。
class Solution:
    def crawl(self, startUrl: str, htmlParser: 'HtmlParser') -> List[str]:
        def get_host(url: str):
            i0 = url.find("//")
            i1 = url.find("/", i0 + 2)
            return url[:i1]

        target = get_host(startUrl)
        queue = deque(htmlParser.getUrls(startUrl))
        res = {startUrl}
        while queue:
            cur = queue.popleft()
            if cur not in res:
                if cur.startswith(target):
                    res.add(cur)
                    queue.extend(htmlParser.getUrls(cur))
        return list(res)

