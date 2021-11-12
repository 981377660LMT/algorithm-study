1. String
   String 是最常用的一种数据类型，普通的 key/value 存储都可以归为此类，这里就不所做解释了。
   String 在 redis 内部存储默认就是一个字符串，被 redisObject 所引用，当遇到 incr,decr 等操作时会转成数值型进行计算，此时 redisObject 的 encoding 字段为 int。
2. List
   Redis list 的实现为一个双向链表，即可以支持反向查找和遍历，更方便操作，不过带来了部分额外的内存开销，Redis 内部的很多实现，包括发送缓冲队列等也都是用的这个数据结构。
   **只允许三个人登录的场景**: lpush、rpush、lpop、rpop
   twitter 的关注列表、粉丝列表等都可以用 Redis 的 list 结构来实现
3. Hash
   Redis Hash 对应 Value 内部实际就是一个 HashMap，实际这里会有 2 种不同实现，这个 Hash 的成员比较少时 Redis 为了节省内存会采用类似**一维数组**的方式来紧凑存储，而不会采用真正的 HashMap 结构
   假设有多个用户及对应的用户信息，可以用来存储以**用户 ID 为 key**，**将用户信息序列化为比如 json 格式做为 value 进行保存。**
4. Set
   set 的内部实现是一个 value 永远为 null 的 HashMap。实际就是通过计算 hash 的方式来快速排重的，这也是 set 能提供判断一个成员是否在集合内的原因
   在微博应用中，**每个用户关注的人存在一个集合中**，就很容易实现求两个人的**共同好友**功能。(交集)
5. Sorted Set
   一个有序集合的每个成员带有分数，用于进行排序。
   Redis sorted set 的内部使用 HashMap 和跳跃表(SkipList)来保证数据的存储和有序
   HashMap 里放的是成员到 score 的映射，而**跳跃表里存放的是所有的成员**，排序依据是 **HashMap 里存的 score**,使用跳跃表的结构可以获得比较高的查找效率，并且在实现上比较简单。
   twitter 的 public timeline 可以以发表时间作为 score 来存储，这样获取时就是自动按时间排好序的。
   又比如用户的**积分排行榜**需求就可以通过有序集合实现
