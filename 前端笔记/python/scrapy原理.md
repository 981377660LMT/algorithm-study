http://kaito-kidd.com/2016/11/01/scrapy-code-analyze-architecture/
Scrapy 是用纯 Python 实现一个为了爬取网站数据、提取结构性数据而编写的应用框架，用途非常广泛，它的设计模式主要是`Template Method Pattern`

1. 五大核心类；Spiders(爬虫类)，Scrapy Engine(引擎),Scheduler(调度器),Downloader(下载器),Item Pipeline(处理管道)。

- Spiders:开发者自定义的一个类，用来解析网页并抓取指定 url 返回的内容，`产生requests。`，处理 response 决定继续还是进管道
- Scrapy Engine:控制整个系统的数据处理流程，并进行事务处理的触发。
- Scheduler：负责管理任务、过滤任务、输出任务的调度器，存储、去重任务都在此控制；(可以被 redis 代替)
- Downloader：抓取网页信息提供给 engine，进而转发至 Spiders。
- Item Pipeline:负责处理 Spiders 类提取之后的数据。比如清理 HTML 数据、验证爬取的数据(检查 item 包含某些字段)、查重(并丢弃)、将爬取结果保存到数据库中

  还有两个模块：
  Downloader middlewares：介于引擎和下载器之间，可以在网页在下载前、后进行逻辑处理；
  Spider middlewares：介于引擎和爬虫之间，在向爬虫输入下载结果前，和爬虫输出请求 / 数据后进行逻辑处理；

2. 9 个步骤
   scrapy 分为 9 个步骤。
   1、Spiders 需要初始的 `start_url` 或则函数 `stsrt_requests`,会在内部生成 Requests 给 Engine；
   2、Engine 将 requests 发送给 Scheduler;
   3、Engine 从 Scheduler 那获取 requests,交给 Download 下载；
   4、在交给 Dowmload 过程中会经过` Downloader Middlewares`(经过 process_request 函数)；
   5、Dowmloader 下载页面后生成一个 response，这个 response 会传给 Engine，这个过程中又经过了 `Downloader Middlerwares`(经过 process_response 函数)，在传送中出错的话经过 process_exception 函数；
   6、Engine 将从 Downloader 那传送过来的 response 发送给 Spiders 处理，这个过程经过 `Spiders Middlerwares`(经过 process_spider_input 函数)；
   7、Spiders 处理这个 response，返回 Requests 或者 Item 两个类型，传给 Engine，这个过程又经过 `Spiders Middlewares`(经过 porcess_spider_output 函数)；
   8、Engine 接收返回的信息，如果使 Item，将它传给 Items Pipeline 中；如果是 Requests,将它传给 Scheduler，继续爬虫；
   9、重复第三步，直至没有任何需要爬取的数据

3. scrapy 入口到运行
   Scrapy 究竟是如何运行起来的。
   scrapy crawl <spider_name>

   我们执行的 scrapy 命令从何而来？
   实际上，当你成功安装好 Scrapy 后，使用如下命令，就能找到这个命令文件，这个文件就是 Scrapy 的运行入口：

   ```
   $ which scrapy
   /usr/local/bin/scrapy
   ```

   使用编辑打开这个文件，你会发现，它其实它就是一个 Python 脚本，而且代码非常少。

   ```Python
   import re
   import sys

   from scrapy.cmdline import execute

   if __name__ == '__main__':
       sys.argv[0] = re.sub(r'(-script\.pyw|\.exe)?$', '', sys.argv[0])
       sys.exit(execute())
   ```

   **execute** 主要工作包括配置初始化、命令解析、爬虫类加载、运行爬虫这几步。

   1. 初始化项目配置：配置是有优先级的

   ```Python
   # default_settings.py

   # 下载器类
   DOWNLOADER = 'scrapy.core.downloader.Downloader'
   # 调度器类
   CHEDULER = 'scrapy.core.scheduler.Scheduler'
   # 调度队列类
   SCHEDULER_DISK_QUEUE = 'scrapy.squeues.PickleLifoDiskQueue'
   SCHEDULER_MEMORY_QUEUE = 'scrapy.squeues.LifoMemoryQueue'
   SCHEDULER_PRIORITY_QUEUE = 'scrapy.pqueues.ScrapyPriorityQueue'
   ```

   有没有感觉比较奇怪，默认配置中配置了这么多类模块，这是为什么？
   这其实是 Scrapy 特性之一，它这么做的好处是：**任何模块都是可替换的。**
   例如，你觉得默认的调度器功能不够用，那么你就可以按照它定义的接口标准，自己实现一个调度器，然后在自己的配置文件中，注册自己的调度器类，那么 Scrapy 运行时就会加载你的调度器执行了，这极大地提高了我们的灵活性！
   所以，**只要在默认配置文件中配置的模块类，都是可替换的。**

   2. 检查运行环境是否在项目中
      运行环境是否在爬虫项目中的依据就是能否找到 scrapy.cfg 文件，如果能找到，则说明是在爬虫项目中，否则就认为是执行的全局命令。
   3. 组装命令实例集合
      scrapy 包括很多命令，例如 scrapy crawl 、 scrapy fetch 等等，那这些命令是从哪来的？答案就在 `_get_commands_dict` 方法中
      这个过程主要是，**导入 commands 文件夹下的所有模块**，`最终生成一个 {cmd_name: cmd_class} 字典集合`，如果用户在配置文件中也配置了自定义的命令类，也会追加进去。也就是说，我们自己也可以编写自己的命令类，然后追加到配置文件中，之后就可以使用自己定义的命令了

   ```Python
   def _iter_command_classes(module_name):
       # 迭代这个包下的所有模块 找到ScrapyCommand的子类
       for module in walk_modules(module_name):
           for obj in vars(module).values():
               if inspect.isclass(obj) and \
                       issubclass(obj, ScrapyCommand) and \
                       obj.__module__ == module.__name__:
                   yield obj
   ```

   4. 解析命令
      加载好命令类后，就开始解析我们具体执行的哪个命令了，解析逻辑比较简单：

   ```Python
   def _pop_command_name(argv):
       i = 0
       for arg in argv[1:]:
           if not arg.startswith('-'):
               del argv[i]
               return arg
           i += 1
   ```

   这个过程就是解析命令行，例如执行 scrapy crawl <spider_name>，这个方法会解析出 crawl，`通过上面生成好的命令类的字典集合，就能找到 commands 目录下的 crawl.py 文件`，最终执行的就是它的 Command 类。

   5. 解析命令行参数
      调用 `cmd.process_options` 方法解析我们的参数：
      固定参数解析交给父类处理，例如输出位置等。其余不同的参数由不同的命令类解析。
   6. 初始化 CrawlerProcess
      一切准备就绪，最后初始化 CrawlerProcess 实例，然后运行对应命令实例的 run 方法。
      我们开始运行一个爬虫一般使用的是 scrapy crawl <spider_name>，也就是说最终调用的是 `commands/crawl.py 的 run 方法`：

   ```Python
   cmd.crawler_process = CrawlerProcess(settings)
   _run_print_help(parser, _run_command, cmd, args, opts)

   def run(self, args, opts):
       if len(args) < 1:
           raise UsageError()
       elif len(args) > 1:
           raise UsageError("running 'scrapy crawl' with more than one spider is no longer supported")
       spname = args[0]

       self.crawler_process.crawl(spname, **opts.spargs)
       self.crawler_process.start()
   ```

   run 方法中调用了 CrawlerProcess 实例的 `crawl 和 start 方法`，就这样整个爬虫程序就会运行起来了。
   我们先来看 CrawlerProcess 初始化：

   ```Python
   class CrawlerProcess(CrawlerRunner):
       def __init__(self, settings=None):
           # 调用父类初始化
           super(CrawlerProcess, self).__init__(settings)
           # 信号和log初始化
           install_shutdown_handlers(self._signal_shutdown)
           configure_logging(self.settings)
           log_scrapy_info(self.settings)

   # 其中，构造方法中调用了父类 CrawlerRunner 的构造方法：
   class CrawlerRunner(object):
       def __init__(self, settings=None):
           if isinstance(settings, dict) or settings is None:
               settings = Settings(settings)
           self.settings = settings
           # 获取爬虫加载器
           self.spider_loader = _get_spider_loader(settings)
           self._crawlers = set()
           self._active = set()

   # 初始化时，调用了 _get_spider_loader方法：
   def _get_spider_loader(settings):
       # 读取配置文件中的SPIDER_MANAGER_CLASS配置项
       if settings.get('SPIDER_MANAGER_CLASS'):
           warnings.warn(
               'SPIDER_MANAGER_CLASS option is deprecated. '
               'Please use SPIDER_LOADER_CLASS.',
               category=ScrapyDeprecationWarning, stacklevel=2
           )
       cls_path = settings.get('SPIDER_MANAGER_CLASS',
                               settings.get('SPIDER_LOADER_CLASS'))
       loader_cls = load_object(cls_path)
       try:
           verifyClass(ISpiderLoader, loader_cls)
       except DoesNotImplement:
           warnings.warn(
               'SPIDER_LOADER_CLASS (previously named SPIDER_MANAGER_CLASS) does '
               'not fully implement scrapy.interfaces.ISpiderLoader interface. '
               'Please add all missing methods to avoid unexpected runtime errors.',
               category=ScrapyDeprecationWarning, stacklevel=2
           )
       return loader_cls.from_settings(settings.frozencopy())

   # 这里会读取默认配置文件中的 spider_loader项，默认配置是 spiderloader.SpiderLoader类，从名字我们也能看出来，这个类是用来加载我们编写好的爬虫类的

   ```

   爬虫加载器会加载所有的爬虫脚本，**最后生成一个 {spider_name: spider_cls} 的字典**，所以我们在执行 scarpy crawl <spider_name> 时，Scrapy 就能找到我们的爬虫类。

   7. 运行爬虫
      CrawlerProcess 初始化完之后，调用它的 **crawl** 方法：

   ```Python
   @defer.inlineCallbacks
   def crawl(self, *args, **kwargs):
       assert not self.crawling, "Crawling already taking place"
       self.crawling = True

       try:
           # 到现在 才是实例化一个爬虫实例
           self.spider = self._create_spider(*args, **kwargs)
           # 创建引擎
           self.engine = self._create_engine()
           # 调用爬虫类的start_requests方法
           start_requests = iter(self.spider.start_requests())
           # 执行引擎的open_spider 并传入爬虫实例和初始请求
           yield self.engine.open_spider(self.spider, start_requests)
           yield defer.maybeDeferred(self.engine.start)
       except Exception:
           if six.PY2:
               exc_info = sys.exc_info()

           self.crawling = False
           if self.engine is not None:
               yield self.engine.close()

           if six.PY2:
               six.reraise(*exc_info)
           raise

   def _create_spider(self, *args, **kwargs):
       return self.spidercls.from_crawler(self, *args, **kwargs)
   ```

   到这里，才会对我们的爬虫类创建一个实例对象，然后创建引擎，之后调用爬虫类的 start_requests 方法获取种子 URL，最后交给引擎执行。

   最后来看 Cralwer 是如何开始运行的额，也就是它的 **start** 方法：

   ```Python
   def start(self, stop_after_crawl=True):
       if stop_after_crawl:
           d = self.join()
           if d.called:
               return
           d.addBoth(self._stop_reactor)
       reactor.installResolver(self._get_dns_resolver())
       # 配置reactor的池子大小(可修改REACTOR_THREADPOOL_MAXSIZE调整)
       tp = reactor.getThreadPool()
       tp.adjustPoolsize(maxthreads=self.settings.getint('REACTOR_THREADPOOL_MAXSIZE'))
       reactor.addSystemEventTrigger('before', 'shutdown', self.stop)
       # 开始执行
       reactor.run(installSignalHandlers=False)
   ```

   在这里有一个叫做 reactor 的模块。reactor 是个什么东西呢？它是 Twisted 模块的事件管理器，我们只要把需要执行的事件注册到 reactor 中，然后调用它的 run 方法，它就会帮我们执行注册好的事件，如果遇到网络 IO 等待，它会自动帮切换到可执行的事件上，非常高效。

   在这里我们不用深究 reactor 是如何工作的，你可以把它想象成一个线程池，只是采用注册回调的方式来执行事件。

4. 核心组件初始化

   1. 爬虫类
      Crawler 的 crawl 方法，我们来看这个方法三步：
      **首先创建了爬虫实例**，然后创建了**引擎**，最后把**爬虫交给引擎来处理了**。
      实例化爬虫比较有意思，它不是通过普通的构造方法进行初始化，而是调用了类方法 from_crawler 进行的初始化

   ```Python

   <!-- scrapy.Spider 类 -->


   @classmethod
   def from_crawler(cls, crawler, *args, **kwargs):
       spider = cls(*args, **kwargs)
       spider._set_crawler(crawler)
       return spider

   def _set_crawler(self, crawler):
       self.crawler = crawler
       # 把settings对象赋给spider实例
       self.settings = crawler.settings
       crawler.signals.connect(self.close, signals.spider_closed)


    class Spider(object_ref):
        name = None
        custom_settings = None

        def __init__(self, name=None, **kwargs):
            # name必填
            if name is not None:
                self.name = name
            elif not getattr(self, 'name', None):
                raise ValueError("%s must have a name" % type(self).__name__)
            self.__dict__.update(kwargs)
            # 如果没有设置start_urls 默认是[]
            if not hasattr(self, 'start_urls'):
                self.start_urls = []
   ```

   这里就是我们平时编写爬虫类时，最常用的几个属性：name、start_urls、custom_settings：

   - name：在运行爬虫时通过它找到我们编写的爬虫类；
   - start_urls：抓取入口，也可以叫做种子 URL；
   - custom_settings：爬虫自定义配置，会覆盖配置文件中的配置项；

   2. 引擎
      Crawler 的 crawl 方法，紧接着就是创建引擎对象
      也就是 `_create_engine` 方法，看看初始化时都发生了什么？

      ```Python
      class ExecutionEngine(object):
        """引擎"""
        def __init__(self, crawler, spider_closed_callback):
            self.crawler = crawler
            # 这里也把settings配置保存到引擎中
            self.settings = crawler.settings
            # 信号
            self.signals = crawler.signals
            # 日志格式
            self.logformatter = crawler.logformatter
            self.slot = None
            self.spider = None
            self.running = False
            self.paused = False
            # 从settings中找到Scheduler调度器，找到Scheduler类
            self.scheduler_cls = load_object(self.settings['SCHEDULER'])
            # 同样，找到Downloader下载器类
            downloader_cls = load_object(self.settings['DOWNLOADER'])
            # 实例化Downloader
            self.downloader = downloader_cls(crawler)
            # 实例化Scraper 它是引擎连接爬虫类的桥梁
            self.scraper = Scraper(crawler)
            self._spider_closed_callback = spider_closed_callback
      ```

      主要是对其他几个核心组件进行定义和初始化，主要包括包括：Scheduler、Downloader、Scrapyer，其中 Scheduler 只进行了类定义，没有实例化。
      引擎是整个 Scrapy 的核心大脑(**Controller**)，它负责管理和调度这些组件，让这些组件更好地协调工作。

   3. **调度器**
      调度器的初始化主要做了 2 件事：

      实例化`请求指纹过滤器`：主要用来过滤重复请求；
      定义不同类型的`任务队列`：优先级任务队列、基于磁盘的任务队列、基于内存的任务队列；
      调度器默认定义了 2 种队列类型：
      基于磁盘的任务队列：在配置文件可配置存储路径，每次执行后会把队列任务保存到磁盘上；
      基于内存的任务队列：每次都在内存中执行，下次启动则消失；

      ```Python
      # 基于磁盘的任务队列(后进先出)
      SCHEDULER_DISK_QUEUE = 'scrapy.squeues.PickleLifoDiskQueue'
      # 基于内存的任务队列(后进先出)
      SCHEDULER_MEMORY_QUEUE = 'scrapy.squeues.LifoMemoryQueue'
      # 优先级队列
      SCHEDULER_PRIORITY_QUEUE = 'queuelib.PriorityQueue'
      ```

      如果我们在配置文件中定义了 JOBDIR 配置项，那么每次执行爬虫时，`都会把任务队列保存在磁盘中，下次启动爬虫时可以重新加载继续执行我们的任务。`
      如果没有定义这个配置项，那么默认使用的是内存队列。
      默认定义的这些队列结构都是后进先出的，什么意思呢？
      如果生成一个抓取任务，放入到任务队列中，那么下次抓取就会从任务队列中先获取到这个任务，优先执行。
      这么实现意味什么呢？其实意味着：**Scrapy 默认的采集规则是深度优先！**
      如何改变这种机制，`变为广度优先采集呢`？这时候我们就要看一下 `scrapy.squeues `模块了，在这里定义了很多种队列：
      如果我们想把抓取任务改为广度优先，我们**只需要在配置文件中把队列类修改为先进先出队列类就可以了**！从这里我们也可以看出，`Scrapy 各个组件之间的耦合性非常低，每个模块都是可自定义的。`

   4. 下载器
      主要是初始化了**下载处理器**、**下载器中间件管理器**以及从配置文件中拿到抓取请求控制的相关参数。
      DownloaderMiddlewareManager 实例化过程：

      ```Python
      class DownloaderMiddlewareManager(MiddlewareManager):
          """下载中间件管理器"""
            component_name = 'downloader middleware'

            @classmethod
            def _get_mwlist_from_settings(cls, settings):
                # 从配置文件DOWNLOADER_MIDDLEWARES_BASE和DOWNLOADER_MIDDLEWARES获得所有下载器中间件
                return build_component_list(
                    settings.getwithbase('DOWNLOADER_MIDDLEWARES'))

            def _add_middleware(self, mw):
                # 定义下载器中间件请求、响应、异常一串方法
                if hasattr(mw, 'process_request'):
                    self.methods['process_request'].append(mw.process_request)
                if hasattr(mw, 'process_response'):
                    self.methods['process_response'].insert(0, mw.process_response)
                if hasattr(mw, 'process_exception'):
                    self.methods['process_exception'].insert(0, mw.process_exception)
      ```

      下载器中间件管理器继承了 MiddlewareManager 类，然后重写了 \_add_middleware 方法，为下载行为定义默认的`下载前、下载后、异常时`对应的处理方法。
      中间件这么做的好处是什么？
      从某个组件流向另一个组件时，会经过一系列中间件，每个中间件都定义了自己的处理流程，相当于一个个管道，输入时可以针对数据进行处理，然后送达到另一个组件，另一个组件处理完逻辑后，又经过这一系列中间件，这些中间件可再针对这个响应结果进行处理，最终输出。

      5. Scraper
         下载器实例化完了之后，回到引擎的初始化方法中，然后就是实例化 Scraper
         **这个类没有在架构图中出现，但这个类其实是处于 Engine、Spiders、Pipeline 之间，是连通这三个组件的桥梁。**
         作用主要是:判断 request 结果是 dict 还是 Response 类以及处理 Spider 中间件，`弄一个这样的辅助类解耦`

         这些组件各司其职，相互协调，共同完成爬虫的抓取任务，而且从代码中我们也能发现，**每个组件类都是定义在配置文件中的，也就是说我们可以实现自己的逻辑，然后替代这些组件，这样的设计模式也非常值得我们学习。**

5. 核心抓取流程

   1. 运行入口
      来看 Cralwer 的 crawl 方法：

      ```Python
        @defer.inlineCallbacks
        def crawl(self, *args, **kwargs):
            assert not self.crawling, "Crawling already taking place"
            self.crawling = True
            try:
                # 创建爬虫实例
                self.spider = self._create_spider(*args, **kwargs)
                # 创建引擎
                self.engine = self._create_engine()
                # 调用spider的start_requests 获取种子URL
                start_requests = iter(self.spider.start_requests())
                # 调用engine的open_spider 交由引擎调度
                yield self.engine.open_spider(self.spider, start_requests)
                yield defer.maybeDeferred(self.engine.start)
            except Exception:
                if six.PY2:
                    exc_info = sys.exc_info()
                self.crawling = False
                if self.engine is not None:
                    yield self.engine.close()
                if six.PY2:
                    six.reraise(*exc_info)
                raise
      ```

      这里首先会创建出爬虫实例，然后创建引擎，之后调用了 spider 的 start_requests 方法，这个方法就是我们平时写的最多爬虫类的父类,它在 `spiders/__init__.py` 中定义

      ```Python
      def start_requests(self):
          # 根据定义好的start_urls属性 生成种子URL对象
          for url in self.start_urls:
              yield self.make_requests_from_url(url)

      def make_requests_from_url(self, url):
          # 构建Request对象
          return Request(url, dont_filter=True)
      ```

   2. 构建请求
      用了 engine 的 open_spider
   3. 引擎调度
   4. 调度器
   5. Scraper
   6. 循环调度
   7. 请求入队
   8. 指纹过滤
   9. 下载请求
   10. 处理下载结果
   11. 回调爬虫
   12. 处理输出
   13. CrawlerSpider

6. 总总结结
   Scrapy 的每个模块的实现都非常纯粹，每个组件都通过`配置文件定义连接起来`，如果想要扩展或替换，只需定义并实现自己的处理逻辑即可，其他模块均不受任何影响，所以我们也可以看到，业界有非常多的 Scrapy 插件，都是通过此机制来实现的。
   虽然它只是个单机版的爬虫框架，但我们可以非常方便地编写插件，或者自定义组件替换默认的功能，从而定制化我们自己的爬虫，最终可以实现一个功能强大的爬虫框架，例如分布式、代理调度、并发控制、可视化、监控等功能，它的灵活度非常高。
