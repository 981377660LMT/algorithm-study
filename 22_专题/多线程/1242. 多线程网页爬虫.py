# 1242. 多线程网页爬虫
# https://leetcode.cn/problems/web-crawler-multithreaded/description/?envType=problem-list-v2&envId=concurrency
#
# 给你一个初始地址 startUrl 和一个 HTML 解析器接口 HtmlParser，请你实现一个 多线程的网页爬虫，用于获取与 startUrl 有 相同主机名 的所有链接。
# 以 任意 顺序返回爬虫获取的路径。
# 爬虫应该遵循：
# - 从 startUrl 开始
# - 调用 HtmlParser.getUrls(url) 从指定网页路径获得的所有路径。
# - 不要抓取相同的链接两次。
# - 仅浏览与 startUrl 相同主机名 的链接。
#
# !读写线程不安全的共享变量需要加锁。
# !按照拓扑序并发.

import os
from threading import Lock
from typing import List
from concurrent.futures import ThreadPoolExecutor
from urllib.parse import urlparse
from collections import deque


class HtmlParser(object):
    def getUrls(self, url: str) -> List[str]: ...


cpu_count = os.cpu_count()


class Solution:
    def crawl(self, startUrl: str, htmlParser: "HtmlParser") -> List[str]:
        target_hostname = urlparse(startUrl).netloc

        visited = set()
        visited.add(startUrl)
        visitedLock = Lock()  # !读写线程不安全的共享变量需要加锁

        def is_same_hostname(url: str) -> bool:
            return urlparse(url).netloc == target_hostname

        def crawl_url(url: str) -> List[str]:
            res = []
            for new_url in htmlParser.getUrls(url):
                if is_same_hostname(new_url):
                    with visitedLock:
                        if new_url in visited:
                            continue
                        visited.add(new_url)
                    res.append(new_url)
            return res

        pool = ThreadPoolExecutor(max_workers=cpu_count * 5)  # type: ignore
        queue = deque([startUrl])

        while queue:
            futures = [pool.submit(crawl_url, url) for url in queue]
            queue.clear()

            for future in futures:
                queue.extend(future.result())

        return list(visited)
