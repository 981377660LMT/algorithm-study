1. 反爬：ip/频率限制/**逆向破解 js**
2. 分布式爬虫
3. **tornado** 天然支持并发
4. retry 装饰器
5. 元类创建单例
6. dict 去重 转 string
7. scrapy 和 scrapy-redis 有什么区别？
   scrapy 是一个 Python 爬虫框架，爬取效率极高，具有高度定制性，但是不支持分布式
   。
   而 scrapy-redis 一套基于 redis 数据库、运行在 scrapy 框架之上的组件，`它是将 scrapy 框架中 Scheduler 替换为 redis 数据库，实现队列管理共享`。可以让 scrapy 支持分布式策略，Slaver 端共享 Master 端 redis 数据库里的 item 队列、请求队列和请求指纹集合。优点：
   1、可以充分利用多台机器的带宽；
   2、可以充分利用多台机器的 IP 地址。
8. 什么是主从同步？
   当客户端向从服务器发送 SLAVEOF 命令，要求从服务器复制主服务器时，从服务器首先需要执行同步操作，也即是，`将从服务器的数据库状态更新至主服务器当前所处的数据库状态`
9. scrapy 框架运行的机制？
   从 start_urls 里获取第一批 url 并发送请求，请求由引擎交给调度器入请求队列，获取完毕后，调度器将请求队列里的请求交给下载器去获取请求对应的响应资源，并将响应交给自己编写的解析方法做提取处理：1. 如果提取出需要的数据，则交给管道文件处理；2. 如果提取出 url，则继续执行之前的步骤（发送 url 请求，并由引擎将请求交给调度器入队列...)，直到请求队列里没有请求，程序结束。
10. 简要介绍下 scrapy 框架及其优势
    scrapy 是一个快速(fast)、高层次(high-level)的基于 Python 的 Web 爬虫构架，用于抓取 Web 站点并从页面中提取结构化的数据。scrapy 使用了 `Twisted` 异步网络库来处理网络通讯。scrapy 框架的异步机制是基于 twisted 异步网络框架处理的，在 settings.py 文件里能够设置具体的并发量数值（`默认是并发量 16`）。

    scrapy 框架的优点： 1)更容易构建大规模的抓取项目; 2)异步处理请求速度非常快; 3)可以使用自动调节机制自动调整爬行速度。
    其 `parse->yield item->pipeline` 流程是全部爬虫的固有模式。
    构造形式主要分 spider.py pipeline.py item.py decorator.py middlewares.py setting.py。

11. scrapy 和 requests 的使用情况？
    requests 是 polling 方式的，会被网络阻塞，不适合爬取大量数据
    scapy 底层是异步框架 twisted ，并发是最大优势

12. Scrapy 的优缺点?
    （1）优势：scrapy 是异步的

    采起可读性更强的 `xpath` 代替正则强大的统计和 log 系统，同时在不一样的 url 上爬行支持 shell 方式，方便独立调试写 `middleware`,方便写一些统一的过滤器，经过管道的方式存入数据库

    （2）缺点：基于 python 的爬虫框架，扩展性比较差

    基于 twisted 框架，运行中的 exception 是不会干掉 reactor，而且`异步框架出错后是不会停掉其余任务的，数据出错后难以察觉。`

13. 三种中间件?
    scrapy 的中间件理论上有三种(`Schduler` Middleware,`Spider` Middleware,`Downloader` Middleware)
    DownloaderMiddleware 主要处理请求 Request 发出去和结果 Response 返回的一些回调，比如说你要加 UserAgent，使用代理，修改 refferer（防盗链），添加 cookie，或者请求异常超时处理啥的
    **常用**：

    1. 爬虫中间件 Spider Middleware
       主要功能是在爬虫运行过程中进行一些处理.

    2. 下载器中间件 Downloader Middleware
       主要功能在请求到网页后,页面被下载时进行一些处理.
       著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

14. 分布式爬虫原理？
    scrapy-redis 实现分布式，其实从原理上来说很简单，这里为描述方便，我们把自己的核心服务器称为 master，而把用于跑爬虫程序的机器称为 slave。
    我们知道，采用 scrapy 框架抓取网页，我们需要首先给定它一些 start_urls，爬虫首先访问 start_urls 里面的 url，再根据我们的具体逻辑，对里面的元素、或者是其他的二级、三级页面进行抓取。而要实现分布式，我们只需要`在这个 starts_urls 里面做文章就行了。`
    我们`在 master 上搭建一个 redis 数据库`（注意这个数据库`只用作 url 的存储`，不关心爬取的具体数据，不要和后面的 mongodb 或者 mysql 混淆），并对每一个需要爬取的网站类型，都开辟一个单独的列表字段。通过设置 `slave 上 scrapy-redis 获取 url 的地址为 master 地址`。这样的结果就是，尽管有多个 slave，然而`大家获取 url 的地方只有一个，那就是服务器 master 上的 redis 数据库`。并且，由于 scrapy-redis 自身的`队列机制，slave 获取的链接不会相互冲突`。这样各个 slave 在`完成抓取任务之后，再把获取的结果汇总到服务器上`（这时的数据存储不再在是 redis，`而是 mongodb 或者 mysql 等存放具体内容的数据库了`）这种方法的还有好处就是程序移植性强，只要处理好路径问题，把 slave 上的程序移植到另一台机器上运行，基本上就是复制粘贴的事情。
15. 如何实现全站数据爬取？

    - 基于手动请求发送+递归解析
    - 基于 CrwalSpider（LinkExtractor，Rule）链接提取器&规则解析器

16. scrapy 去重所用的几种机制,基于内存
    1. dupefilters.py 去重器:需要将 dont_filter 设置为 False 开启去重，默认是 False；对于每一个 url 的请求，**调度器**都会根据请求的相关信息加密得到一个**指纹信息**，并且将指纹信息和 set()集合中得指纹信息进行比对，如果 set()集合中已经存在这个数据，就不在将这个 Request 放入队列中。如果 set()集合中没有，就将这个 Request 对象放入队列中，等待被调度。
       指纹信息:finger_print 包括:`url/method/data` 三个部分
    2. redis 基于内存 更加快捷、速度快、易于管理
    3. 布隆过滤器
17. 爬取数据后使用哪个数据库存储数据的，为什么
18. 如何提升 scrapy 的爬取效率
    (1)增加并发

默认 scrapy 开启的并发线程为 32 个, 可以适当进行增加. 在 settings 配置文件中修改``CONCURRENT_REQUESTS` = 100 值为 100, 并发设置成了为 100.

(2)降低日志级别

在运行 scrapy 时, 会有大量日志信息的输出, 为了减少 CPU 的使用率. 可以设置 log 输出信息为 INFO 或者 ERROR. 在配置文件中编写: `LOG_LEVEL` = ‘INFO’

(3)禁止 cookie

如果不是真的需要 cookie, 则在 scrapy 爬取数据时可以禁止 cookie 从而减少 CPU 的使用率, 提升爬取效率. 在配置文件中编写: `COOKIES_ENABLED` = False.

(4)禁止重试

对失败的 HTTP 进行重新请求(重试)会减慢爬取速度, 因此可以禁止重试. 在配置文件中编写: `RETRY_ENABLED` = False

(5)减少下载超时

如果对一个非常慢的链接进行爬取, 减少下载超时可以能让卡住的链接快速被放弃, 从而提升效率. 在配置文件中进行编写: `DOWNLOAD_TIMEOUT` = 10 超时时间为 10s.
