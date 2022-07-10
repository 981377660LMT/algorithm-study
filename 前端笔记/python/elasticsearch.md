https://juejin.cn/post/6958408979235995655/

https://jishuin.proginn.com/p/763bfbd31fed

1. 我们可以在 Elasticsearch 中执行搜索的各种可能方式有哪些？
   方式一：基于 DSL 检索（最常用） Elasticsearch 提供基于 JSON 的完整查询 DSL 来定义查询。
2. Elasticsearch 是一个分布式的 RESTful 风格的搜索和数据分析引擎。
3. elasticsearch 的倒排索引是什么
   从单词到文档 id 的映射
4. ES 中数据量有多大？查询耗时多少？耗时的瓶颈在哪，怎么去排查的？
5. 对 ES 做了哪些查询的优化？
6. 你这个部分的 ES 索引是怎么设计的？为什么要这用这种数据类型？
7. 分词器是用的哪个？为什么要用这个？
   kuromoji 分词和 ik 分词(ik-smart 以及)

主要是分词功能
将这些文件导入 ElasticSearch，我们还需要把 json 转换为对应的 bulk api 的格式，这篇文章写的非常详细：https://zhuanlan.zhihu.com/p/146511102
运行 es 和 kibana ，导入 json 文件到 Elasticsearch？

如果需要用双引号""精准查询，就用 `bool 查询`里的 must，并指定为 `phrase 类型`搜索；如果需要部分匹配查询，就用 `bool 查询里的 should`，并指定 fuzziness 为 auto。

此外，我们搜索时必须使用分词。`中文分词为 ik 分词，日语分词为 kuromoji 分词`。

中文部分：存储数据时，需要用 ik_max_word 尽可能多的分词，搜索数据时，需要用 ik_smart 尽可能智能地去匹配。

日文部分：存储数据和搜索数据都是 kuromoji 分词。kuromoji 分词的详细介绍：https://juejin.cn/post/6844903854337687559

谷歌浏览器禁止携带 cookie 跨域问题

开发环境下，在发送验证码环节，我使用的是 cookie 和 session 来存储用户的验证码。结果后端无法将发送的验证码保存在 session 里。经过与热心网友的长时间讨论，发现 response 里面根本没有 set-cookie 这个字段！试了一下 Microsoft Edge 浏览器，发现完全没有问题。后来才了解到这是 Chrome 浏览器 80 版本之后的为了防止跨站请求伪造(csrf)而采用的 cookie 策略，`http 默认 Samesite=lax`,无法携带 cookie 进行`跨站`请求;`https 才能开启 Samesite=None`！
注意跨站，不能跨站自然就不能跨域
注意后端要配合 CORS 设置 `ACCESS-CONTROL-ALLOW-CREDENTIALS` 和 ACCESS-CONTROL-ALLOW-ORIGIN

域名在 NameCheap 购买，\*.me 域名可以一年免费；使用 Cloudflare 的服务进行域名解析与 CDN 加速，还可以开启 https 访问；DigitalOcean 购买服务器，使用 github 的学生包可以获得$50 的优惠，并且首次注册送期限两个月的$100 消费额度；1G 内存的服务器每个月$5，4G 内存的服务器每个月$20，支持 Paypal 支付。选择服务器(droplets)前，记得测试一下哪个地方的服务器访问最快(测试地址：http://speedtest-sgp1.digitalocean.com/)。
