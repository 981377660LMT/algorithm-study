[游戏 AI 浅谈](https://blog.csdn.net/jk_chen_acmer/article/details/119904413)

17. 游戏中的人工智能 2 – Advanced Artificial Intelligence

- 层次任务网络(hierarchical tasks network, HTN)
  ![alt text](image-79.png)
  - 世界的状态+行为的条件+行为的影响
  - 任务就是行为的抽象，每个任务都可以有条件和影响
  - 规划(Planning)
    带着目的进行预测
  - Replan
- goal-oriented action planning(GOAP)
  一种基于`规划`的 AI 技术，和前面介绍过的方法相比 GOAP 一般会更适合动态的环境。
  - 目标 + 行动 + 规划器
  - 反向规划(backward planning)。整个规划问题就等价于在有向图上的最短路径问题。
- 蒙特卡洛树搜索(Monte Carlo tree search, MCTS)
  https://ouuan.github.io/post/monte-carlo-tree-search/

  1. State and Action
     ![alt text](image-80.png)
     ![alt text](image-81.png)
  2. Simulation
     AI 利用当前的策略(default policy)快速地完成整个游戏过程。
     ![alt text](image-82.png)
  3. Backpropagate
     到了叶节点后，会`向上更新这次搜索所经过的路径上的每个节点的被访问次数和胜利次数。`
     蒙特卡洛树搜索基于这样一个设定：每次搜索其实就是在模拟玩家的选择，搜索时某个子节点的被访问次数更多，实际游戏中选择这个子节点就更优；而搜索次数越多，对玩家最优选择的模拟就越准确。这样的话，当搜索次数足够多时，每次选择都是对于当前节点的先手玩家而言最优的，就收敛到了 Min-Max 搜索。
  4. Iteration Steps
     ![alt text](image-84.png)
     - 选择（Selection）：在对节点进行选择时，MCTS 会优先选择可拓展的节点。
       在进行拓展时往往还要权衡一些 exploitation 和 exploration，因此我们可以把 UCB 可以作为一种拓展的准则。
       ` 让游戏树向最优的方向扩展，这是蒙特卡洛树搜索的精要所在。`
       ![alt text](image-83.png)
       Exploitaion：开发，选择已知最优的节点。
       Exploration：探索，选择未知的节点。
       有点像模拟退火的思想。
     - 扩展（Expansion）：除非任意一方的输赢使得游戏在 L 结束，否则创建一个或多个子节点并选取其中一个节点 C。
     - 仿真（Simulation）：再从节点 C 开始，用随机策略进行游戏，又称为 playout 或者 rollout。
     - 反向传播（Backpropagation）：使用随机游戏的结果，更新从 C 到 R 的路径上的节点信息。

  MCTS 最主要的优点有两点：

  1. `如果你完全不会一个游戏，只知道它的规则，也可以使用 MCTS。`
     而 Min-Max 搜索必须有一个估价函数。
     当然，如果完全使用原始的基于随机选择的 MCTS，棋力不一定足够高。
  2. 由于对随机采样的利用，可以胜任较大的搜索空间。

- Machine Learning Basic
  机器学习大致可以分为监督学习、无监督学习、半监督学习以及强化学习等几类。
  `监督学习(Supervided Learning，SL)`：data + label，例如分类、回归
  无监督学习：data，例如聚类
  半监督学习：data + label + data，例如半监督分类

  `强化学习(reinforcement learning, RL)`是游戏 AI 技术的基础。在强化学习中我们希望 AI 能够通过和环境的不断互动来学习到一个合理的策略。
  强化学习：agent + environment + action + reward，例如 AlphaGo

  Markov Decision Process(MDP)：状态、动作、奖励、策略、价值函数

- Build Advanced Game AI
  尽管目前基于机器学习的游戏 AI 技术大多还处于试验阶段，但已经有一些很优秀的项目值得借鉴和学习，包括 DeepMind 的 AlphaStar 以及 OpenAI 的 Five 等。
  基于`深度强化学习(deep reinforcement learning, DRL)`的游戏 AI 都是使用一个`深度神经网络`来进行决策，整个框架包括接收游戏环境的观测，利用神经网络获得行为，以及从游戏环境中得到反馈。

  1. State
     以 AlphaStar 为例，智能体可以直接从游戏环境获得的信息包括地图、统计数据、场景中的单位以及资源数据等。
  2. Action
     ![alt text](image-85.png)
  3. Reward
     奖励函数的设计对于模型的训练以及最终的性能都有着重要的影响。在 AlphaStar 中使用了非常简单的奖励设计，智能体仅在获胜时获得+1 的奖励；而在 OpenAI Five 中则采用了更加复杂的奖励函数并以此来鼓励 AI 的进攻性。
  4. Network
     在 AlphaStar 中使用了不同种类的神经网络来处理不同类型的输入数据，比如说对于定长的输入使用了 MLP，对于图像数据使用了 CNN，对于非定长的序列使用了 Transformer，而对于整个决策过程还使用了 LSTM 进行处理。

     MLP：多层感知机
     CNN：卷积神经网络，例如 ResNet
     Transformer：自注意力机制
     LSTM：长短期记忆网络

  5. Training Strategy
     在 AlphaStar 的训练过程中首先使用了监督学习的方式来`从人类玩家的录像中进行学习。`接着，AlphaStar 使用了强化学习的方法来进行自我训练。
     ![三手互搏](image-86.png)
     试验结果分析表明基于监督学习训练的游戏 AI 其行为会比较接近于人类玩家，但基本无法超过人类玩家的水平；而基于强化学习训练的 AI 则可能会有超过玩家的游戏水平，不过需要注意的是使用强化学习可能需要非常多的训练资源。

     **因此对于游戏 AI 到底是使用监督学习还是使用强化学习进行训练需要结合实际的游戏环境进行考虑。对于奖励比较密集的环境可以直接使用强化学习进行训练，而对于奖励比较稀疏的环境则推荐使用监督学习。**

     奖励足够密集：打方块
     奖励不够密集：超级马里奥

  训练在 gpu，运行在 cpu

2.  网络游戏的架构 1 – 基础

- Network Protocols
  人们提出了中间层(intermediate layer)的概念来隔绝掉应用和硬件，使得开发者可以专注于程序本身而不是具体的通信过程。
  ![alt text](image-87.png)
  在现代计算机网络中人们设计了 OSI 模型(OSI model)来对通信过程进行封装和抽象。
  ![alt text](image-88.png)

  1. Socket
  2. TCP
  3. UDP
     除了 TCP 之外人们还开发出了 UDP 这样的轻量级网络协议。UDP 的本质是一个端到端的网络协议，它不需要建立长时间的连接，也不要求发送数据的顺序，因此 UDP 要比 TCP 简单得多。
     ![alt text](image-89.png)

     `对于实时性要求比较高的游戏会优先选择 UDP，而策略类的游戏则会考虑使用 TCP。`在大型网络游戏中还可能会使用复合类型的协议来支持游戏中不同系统的通信需求。

  4. Reliable UDP
     同时现代网络游戏中往往还会对网络协议进行定制。以 TCP 为例，虽然 TCP 协议比较稳定但是效率过于低了，而且网络游戏中出现一定的丢包是可以接受的；而对于 UDP 来说它虽然非常高效但是却不够稳定。
     因此现代网络游戏中往往**会基于 UDP 来定制一个网络协议，这样既可以利用 UDP 的高效性又可以保证数据通信的有序性。**
     ACK(确认消息) 及其相关技术是保证数据可靠通信的基本方法。

     - **ARQ(automatic repeat request，自动重传请求)**是基于 ACK 的错误控制方法，所有的通信算法都要事项 ARQ 的功能。
       滑动窗口协议(sliding window protocol)是一种经典的 ARQ 实现方法，它在发送数据时每次发送窗口大小的包然后检验回复的 ACK 来判断是否出现丢包的情况。（滑动窗口哈希）
     - **Forward Error Correction(FEC，前向纠错)**
       前向纠错码，是一种通过增加冗余数据来实现数据纠错的方法。
       目前常用的 FEC 算法包括异或校验位以及 Reed-Solomon codes 两大类。
       海明码(Hamming code)是一种最简单的纠错码，它通过增加冗余位来实现数据的纠错。
       Reed-Solomon codes 是经典的信息传输算法，它利用 Vandemode 矩阵及其逆阵来恢复丢失的数据。

  `总结一下，在自定义 UDP 时需要考虑 ARQ 和 FEC 两类问题。`

- Clock Synchronization
  有了网络协议后就可以开始对网络游戏进行开发了，不过在具体设计游戏前我们还需要考虑`不同玩家之间的时钟同步(clock synchronization)问题。`
  1. Round Trip Time(RTT)
     RTT 是指从发送数据到接收到回复的时间，它是网络通信中的一个重要指标。
     这个间隔的时间称为 round-trip time(RTT)。`RTT 的概念类似于 ping，不过它们的区别在于 ping 更加偏向于底层而 RTT 则位于顶部的应用层。`
  2. 网络时间协议 NTP（Network Time Protocol）
     NTP 同步原理：
     https://info.support.huawei.com/info-finder/encyclopedia/zh/NTP.html
     只需要从客户端发送请求然后从服务器接收一个时刻就好，这样就可以得到 `4 个时间戳`。如果我们进一步假定网络`上行和下行的延迟是一致的`，我们可以直接计算出 RTT 的时间长短以及两个设备之间的时间偏差。
     实际上我们可以证明在不可靠的通信中是无法严格校准时间的。不过在实践中我们可以通过不断的使用 NTP 算法来得到一系列 RTT 值，然后`把高于平均值 50%的部分丢弃，剩下的 RTT 平均值的 1.5 倍就可以作为真实 RTT 的估计。`
- Remote Procedure Call (RPC 协议)
  尽管利用 socket 我们可以实现客户端和服务器的通信，但`对于网络游戏来说完全基于 socket 的通信是非常复杂的。这主要是因为网络游戏中客户端需要向服务器发送大量不同类型的消息，同时客户端也需要解析相应类型的反馈，这就会导致游戏逻辑变得无比复杂。`
  在现代网络游戏中一般会使用 RPC(remote procedure call)的方式来实现客户端和服务器的通信。`基于 RPC 的技术在客户端可以像本地调用函数的方式来向服务器发送请求，这样使得开发人员可以专注于游戏逻辑而不是具体底层的实现。`
  ![alt text](image-90.png)

  golang 的 rpc
  ![alt text](image-91.png)

  在 RPC 中会大量使用 `IDL(interface definiton language，接口描述语言)`来定义不同的消息形式。
  ![alt text](image-92.png)

  然后在启动时通过 RPC stubs 来通知客户端有哪些 RPC 是可以进行调用的。
  ![alt text](image-93.png)
  stub 类似 python 的存根(.后缀为.pyi 的 stub 存根文件)，是类型定义

- Network Topology
  在设计网络游戏时还需要考虑网络自身的架构。

1. Peer-to-Peer
   最经典的网络架构是 P2P(peer-to-peer)，此时每个客户端之间会直接建立通信。
   很多早期经典的游戏都是使用这样的网络架构来实现联网功能。
   当 P2P 需要集中所有玩家的信息时则可以选择其中一个客户端作为主机，这样其它的客户端可以通过连接主机的方式来实现联机。
   ![alt text](image-94.png)
2. Dedicated Server
   ![alt text](image-95.png)
   从实践结果来看，对于小型的网络游戏 P2P 是一个足够好的架构，而对于大型的商业网络游戏则必须使用 dedicated server 这样的形式。

- **Game Synchronization**
  而在网络游戏中，除了单机游戏都需要的分层外我们还需要考虑不同玩家之间的同步。在理想情况下我们希望客户端只负责处理玩家的输入，整个游戏逻辑都放在服务器端。
  ![alt text](image-96.png)

  1. Snapshot Sync (快照同步)
     代表：MineCraft
     快照同步(snapshot synchronization)是一种相对古老的同步技术。在快照同步中`客户端只负责向服务器发送当前玩家的数据，由服务器完成整个游戏世界的运行。然后服务器会为游戏世界生成一张快照，再发送给每个客户端来给玩家反馈。`
     快照同步可以`严格保证每个玩家的状态都是准确的`，但其缺陷在于它给服务器提出了非常巨大的挑战。因此在实际游戏中一般会降低服务器上游戏运行的帧率来平衡带宽，然后在客户端上通过插值的方式来获得高帧率。
     由于每次生成快照的成本是相对较高的，为了压缩数据我们可以使用`状态的变化量来对游戏状态进行表示`。
     快照同步非常简单也易于实现，`但它基本浪费掉了客户端上的算力同时在服务器上会产生过大的压力。因此在现代网络游戏中基本不会使用快照同步的方式。`
  2. Lockstep Sync (帧同步)
     代表：英雄联盟
     帧同步(lockstep synchronization)是现代网络游戏中非常常用的同步技术。`不同于快照同步完全通过服务器来运行游戏世界，在帧同步中服务器更多地是完成数据的分发工作。`玩家的操作通过客户端发送到服务器上，经过服务器汇总后将当前游戏世界的状态返还给客户端，然后在每个客户端上运行游戏世界。
     ![alt text](image-97.png)

     > Same Input + Same Execution Process = Same State

     - Lockstep Initialization
       使用帧同步时首先需要进行初始化，将客户端上所有的游戏数据与服务器进行同步。这一过程一般是`在游戏 loading 阶段`来实现的。
     - Deterministic Lockstep
       早期的联网游戏问题：`当某个玩家的数据滞后了所有玩家都必须要进行等待。`
       waiting for players...
     - Bucket Synchronization
       为了克服这样的问题，人们提出了 bucket synchronization 这样的策略。
       此时`服务器只会等待 bucket 长度的时间，如果超时没有收到客户端发来的数据就越过去`，看下一个 bucket 时间段能否接收到。通过这样的方式其它玩家就无需一直等待了。
       bucket synchronization`本质是对玩家数据的一致性以及游戏体验进行的一种权衡(trade-off)。`
     - Deterministic(确定性) Difficulties
       帧同步的一大难点在于它要保证不同客户端上游戏世界在相同输入的情况下有着完全一致的输出。
       为了保证输出的确定性我们首先要保证浮点数在不同客户端上的一致性，这可以使用 IEEE 754 标准来实现。
       其次在不同的设备上我们需要保证相关的数学运算函数有一致的行为，对于这种问题则可以使用`查表`的方式来避免实际的计算。
       还要考虑随机数的问题，我们要求`随机数在不同的客户端上也必须是完全一致的`。因此在游戏客户端和服务器进行同步时需要`将随机数种子以及随机数生成算法进行同步`

     - Tracing and Debugging
       对于服务器来说检测客户端发送的数据是否存在 bug 就非常重要。一般来说我们会要求`客户端每隔一段时间就上传本地的 log，由服务器来检查上传数据是否存在 bug。`
       server 比较不同 client 的 checksum
     - Lag and Delay
       为了处理网络延迟的问题我们还可以在客户端上`缓存若干帧`，当然缓存的大小会在一定程度上影响玩家的游戏体验。
       eg: jitter buffer 抖动缓冲器
       另一方面我们还可以把游戏逻辑帧和渲染帧进行分离，然后通过插值的方式来获得更加平滑的渲染效果。
     - Reconnection Problem
       断线重连的机制
       实际上在进行帧同步时每个若干帧会**设置一个关键帧。在关键帧进行同步时还会更新游戏世界的快照，这样可保证即使游戏崩溃了也可以从快照中恢复。**
       为了实现这样的功能可以使用 `quick catch up` 技术，此时我们暂停游戏的渲染把所有的`计算资源用来执行游戏逻辑`，以追上游戏进度。可以理解为渲染 silent 了。
       在服务器端也可以使用类似的技术，从而**帮助掉线的玩家快速恢复到游戏的当前状态**。`实际上网络游戏的观战和回放功能也是使用这样的技术来实现的。`
       ![alt text](image-98.png)
       ![alt text](image-99.png)

       观战：本质是重连。Watching is similar to reconnecting after a client crash.
       回放(replay)：`关键帧 + commands`

     - Lockstep Cheating Issues
       对于帧同步的游戏，玩家可以通过发送虚假的状态来实现作弊行为，这就要求我们实现一些反作弊机制。
       - 多人 PVP：投票机制
       - 双人：服务器记录 checksum。不过这种一般都是 p2p 架构，可以换状态同步。
       - 需要解决的问题：客户端外挂，例如透视挂
     - Lockstep Summary
       总结一下，帧同步会占用更少的带宽也比较适合各种需要实时反馈的游戏。而帧同步的难点主要集中在如何保证在不同客户端上游戏运行的一致性，如何设计断线重连机制等。

  3. State Sync (状态同步)
     代表：Counter-Strike(反恐精英)、守望先锋
     状态同步(state synchronization)是目前大型网游非常流行的同步技术，它的基本思想是把玩家的状态和事件进行同步。
     ![alt text](image-100.png)
     进行状态同步时由客户端提交玩家的状态数据，而服务器则会在收集到所有玩家的数据后运行游戏逻辑，然后把下一时刻的状态分发给所有的客户端。

     server 端有一个完整的世界，只会把部分信息发给客户端

     - Authorized and Replicated Clients
       状态同步中服务器称为 authorized server，它是整个游戏世界的绝对权威；(1P)
       而玩家的本地客户端称为 authorized client，它是玩家操作游戏角色的接口；
       在`其他玩家视角下的本地客户端`则称为 replicated client，表示它们仅仅是 authorized client 的一个副本。(3P)
       ![alt text](image-101.png)
     - State Synchronization Example
       authorized client 执行了某种行为时首先会向服务发送相关的数据，**服务器进行计算再发布给所有的客户端**，驱动 replicated client 执行 authorized client 的行为。（活在缸中之脑，每个人活在自己的世界）
       这样的好处在于我们`无需要求每个客户端上的模拟是严格一致的，整个游戏世界本质上仍然是由统一的服务器进行驱动。`
       因为是服务器计算的，所以不存在歧义。
     - Dumb Client Problem
       由于游戏角色的所有行为都需要经过服务器的确认才能执行，状态同步会产生 dumb client 的问题，即`玩家视角下角色的行为可能是滞后的。`
       要缓解这样的问题可以`在客户端上对玩家的行为进行预测(client-side prediction)`。比如说当角色需要进行移动时首先在本地移动半步，然后等服务器传来确定的消息后再进行对齐，这样就可以改善玩家的游戏体验。在守望先锋中就使用了这样的方式来保证玩家顺畅的游玩。
       ![alt text](image-102.png)

       > 有点像 preload(预加载)，先斩后奏

       - Server Reconciliation
         由于网络波动的存在，来自服务器的确认消息往往会滞后于本地的预测。因此我们可以`使用一个 buffer 来记录游戏角色的状态`，这样当收到服务器的消息时首先跟 buffer 中的状态进行检验。当 buffer 中的状态和服务器的数据不一致时就需要根据服务器的数据来矫正玩家状态。

         例子：移动被服务器打回

         当然这样的机制对于网络条件不好的玩家是不太公平的，他们的角色状态会不断地被服务器修正。

     - Packet Loss
       对于丢包的问题在服务器端也会维护一个小的 buffer 来储存玩家的状态。如果 `buffer 被清空则说明可能出现了掉线的情况`，此时服务器会复制玩家上一个输入来维持游戏的运行。

     帧同步和状态同步两种主流同步技术的对比如下：
     ![alt text](image-103.png)

4. 网络游戏的架构 2 – 进阶

- 角色位移同步(character movement replication)
  网络环境的不稳定，玩家操作角色在自己视角和其他玩家视角下的行为往往会有一定的延迟，即角色在其他玩家视角下的动作会滞后于操作玩家的第一视角。
  在这种情况下我们可以使用内插(interpolation)和外插(extrapolation)两种插值方法来缓解延迟。

  - 内插是指利用已知的控制点来获得中间的状态。当网络存在波动时利用内插的方法可以保证角色的动作仍然是足够平滑的。
  - 外插的本质是利用已有的信息来预测未来的状态
    **dead reckoning 算法，航位推测法**

    回到游戏领域的应用中来，projective velocity blending (PVB)是一种使用外插来更新角色位置的算法。

    外插在处理碰撞时会容易产生严重的物体穿插问题。要处理这样的情况`一般需要切换到本地物理引擎来处理碰撞问题，比如说在看门狗 2 中就使用了这样的方法。`

  简单总结一下，对于玩家操作角色`经常出现瞬移或是具有很大加速度的情况比较适合内插(例如MOBA)，对于操作角色比较符合物理规律的情况比较适合外插(例如载具)`，而在一些大型在线游戏中还会同时结合这两种插值方法来提升玩家的游戏体验。

- 命中判定(hit registration)
  以 FPS 游戏为例，从玩家开枪到击中敌人这一过程实际上有着非常大的一段延迟，在各种不确定因素的影响下如何判断玩家确实击中目标就需要一些专门的设计。

  1. Challenges：

     - 确定敌人在什么位置。由于网络延迟和插值算法的存在，玩家视角下的目标是落后于服务器上目标的真实位置的。
     - 如何判定是否击中了目标。

  2. Client-Side Hit Detection
     命中判定的目标是保证游戏中的**玩家就是否命中的问题能够达成一个共识(consensus)**。目前主流的处理方法包括`在客户端上进行检测，称为 client-side hit detection，以及在服务器端进行判断的 server-side hit registration。`

     在客户端上进行检测时的基本思想是一切以玩家客户端视角下的结果为准。玩家开枪后的弹道轨迹以及击中判断都先在本地进行，然后发送到服务器上再进行验证。
     在服务器上会对玩家的行为进行一些验证(甚至是猜)从而保证确实击中了目标。当然在实际的游戏中这个验证过程是相当复杂的，涉及到大量的验证和反作弊检测。
     在客户端上进行命中检测的优势在于它`非常高效`而且可以减轻服务器的负担，但它的核心问题在于它`不够安全`，一旦客户端被破解或是网络消息被劫持就需要非常复杂的反作弊系统来维持游戏平衡。

  3. Server-Side Hit Registration
     角色的位置和状态是领先于玩家视角的，当玩家开枪时目标很可能已经移动到其它的位置上了。从这样的角度来看，玩家很难命中移动中的目标。
     因此我们需要对网络延迟进行一定的`延迟补偿(lag compensation)`。当服务器收到射击的消息时不会直接使用当前的游戏状态，而是根据延迟使用之前保存的游戏状态。
  4. Cover Problems
     对于掩体的问题，由于网络延迟的存在可能会出现射击者优势或窥视者优势的情况。
     为了缓解延迟的问题在游戏设计时还可以`利用动作前摇(startup frames)来掩盖掉网络同步。`
     类似地，也可以使用各种特效来为服务器同步争取时间。

- 大型多人在线游戏(massively multiplayer online game, MMOG) Network Architecture
  刀剑神域!

  MMOG 中一般会有各种子系统来组成整体的玩法系统。

  1. Link Layer
     连接层(link layer)是玩家和游戏建立连接的一层。在 MMO 中为了保护服务器不受攻击，我们需要`单独的连接层来分离玩家和服务器数据。`相当于 GateWay。
  2. Lobby
  3. Character Server
     由于 MMO 中玩家的数据往往非常巨大，玩家的数据一般会保存在一个专门的服务器上，称为 character server。
  4. Trading System
     在设计交易系统时需要保证系统有足够高的安全性。
  5. Social System
  6. Matchmaking
  7. Data Storage
     - Relational Data Storage
       玩家数据和游戏数据都会使用关系数据库进行存储。现代网络游戏中往往还会结合分布式的技术进行存储。
     - Non-Relational Data Storage
       游戏中的日志还有各种 game state 就比较适合使用非关系数据库
     - In-Memory Data Storage
       排行榜，配对。
  8. Distributed System
     - 负载均衡(load balancing)
       有状态服务的`一致性`哈希
     - Servers Management
       如何管理大量同时运行的服务。
       我们可以使用 Apache 或是 etcd 这样的工具来监视和管理各种服务。

- 带宽优化(bandwidth optimization)
  1. 数据压缩(data compression)
     在网络游戏中我们可以把浮点数转换为低精度的定点数来减少数据量。我们甚至可以对游戏地图进行`分区`然后在小区域中使用定点数来缓解低精度数值带来的影响。
  2. Object Relevance
     只同步和玩家相关的对象
     我们可以把整个游戏世界划分为若干个区域，这样每个玩家都会位于某个区域中。不同区域的数据是相互隔绝的，因此在同步时只需要同步区域中的数据即可。
     在一些开放世界游戏中我们不希望出现区域的划分，此时则可以利用 `area of inerest (AOI)`的概念。AOI 的意义在于我们只需要关注玩家附近的情况而无需考虑更远区域的信息。
     ds: 圆、网格划分(带有 notify 机制)、十字链表、PVS(Potential Visibility Set)
  3. Varying Update Frequency
     调整更新频率的方式来降低带宽(`降级`)
     一种经典的策略是根据物体和玩家之间的距离来设置更新频率，**使得距离玩家远的物体更新得慢一些，距离近的更新得快一些。**
- 反作弊(anti-cheat)

  1. Obfuscating Memory(内存混淆)
     修改本地的内存。对于很多在客户端进行校验的游戏只需要修改游戏内存中的数据就可以进行作弊。
     为了应对这种作弊方式，我们可以`对客户端套一个壳来防止侵入`。类似地，也可以对内存数据进行加密来防止侵入。

  2. Verifying Local Files by Hashing
     另一种常见的作弊方法是修改本地的资源文件。我们可以通过对比本地和服务器上资源的 hashing 来进行处理。
  3. Packet Interception and Manipulation
     消息的劫持和篡改是网络游戏中经常遇到的作弊方式，因此对于客户端和服务器发出的包就必须进行加密。
  4. System Software Invoke
     作弊者还可以通过修改底层游戏引擎代码来进行作弊。
     针对这种情况可以`使用各种安全软件来检测游戏引擎是否存在注入`的情况。
     vac，easy(小蓝熊)
  5. AI Cheat
     Detecting Known Cheat Program

- Build a Scalable World
  如何构建一个开放世界
  三种模型：
  1. Zoning(分块)
  2. Instancing
  3. Replication(分层)
     ![alt text](image-104.png)

3. 前沿介绍 1 – Data Oriented Programming，Job System

- Basics of Parallel Programming (并行编程基础)
  进程数据隔离，线程数据共享

  1. Types of Multitasking
     对于多核的计算机我们希望能够充分利用不同的计算核心来提升程序的性能。根据处理器管理任务的不同可以把进程调度分为两种：`抢占式(preemptive multitasking)和非抢占式(non-preemptive multitasking)式。`preemptive multitasking 是由**调度器**来控制任务的切换，而 non-preemptive multitasking 则是由**任务自身**来进行控制。

     调度器：scheduler + interrupt
     任务自身：yield
     ![alt text](image-106.png)

  2. Thread Context Switch
     线程切换耗时：2000 cycles
  3. Parallel Problems in Parallel Computing
     并行案例.
  4. `Data Race` in Parallel Programming
     - lock
     - 原子操作(atomic operation)

- Parallel Framework of Game Engine
  游戏引擎的并行框架
  1. Fixed Multi-Thread
     经典，引擎中的每个系统都各自拥有一个线程。在每一帧开始时会通过线程间的通信来交换数据，然后各自执行自己的任务。
     缺陷在于它很难保证不同线程上负载是一致的。
  2. Thread Fork-Join
     申请一系列线程，当需要执行计算时通过 fork 操作把`不同的计算任务分配到各个线程中并最后汇总到一起。`
     缺陷在于有很多的任务是无法事先预测具体的负载的。
  3. Task Graph
     根据不同任务之间的`依赖性`来决定具体的执行顺序以及需要进行并行的任务
- Job System

  1. Coroutine
     ![alt text](image-107.png)
     有栈协程：使用栈来保存函数切换时的状态。例如 yield。
     无栈协程：不保存函数切换时的状态，当协程切换回来后按照当前的状态继续执行程序。
     `在实践中一般推荐使用基于栈来实现的协程。尽管它在进行切换时的开销要稍微多一些，但可以避免状态改变导致的各种问题。`
  2. Fiber-Based Job System
     基于协程的思想可以实现 fiber-based job system。在这种任务系统中 `job 会通过 fiber 来进行执行`，在线程进行计算时通过 fiber 的切换来减少线程调度的开销

     eg：GMP 模型

     `调度器、job、线程/fiber`

     **在执行计算时根据程序的需要生成大量的 job，然后调度器根据线程负载分配到合适的线程以及线程上的 fiber 中。**
     当 job 出现`依赖`时会把当前的 job 移动到等待区然后执行线程中的下一个 job。这样的方式可以减少 CPU 的等待，提高程序效率。如果出现了线程闲置的情况，调度器会把其他线程中的 job 移动到`闲置线程`中进行计算(偷)。

- 编程范式(programming paradigms)
  1. POP
  2. OOP
     OOP 的问题：
     - `二义性`：角色的攻击行为既可以写在角色身上，也可以写在被攻击者身上。
     - `臃肿`：大量的继承关系，有时很难去查询对象的方法具体是在那个类中实现的。**使用基类的代价。**
     - `性能`：数据往往分布在不同的内存区域上
     - `可测试性`
       要去测试对象的某个方法是否按照我们的期望运作时，`往往需要从头创建出整个对象`，这与单元测试的思想是相违背的。
  3. DOP(Data-Oriented Programming)
     - Cache
     - Locality
     - SIMD
     - LRU
     - Cache Line
- Data-Oriented Programming
  DOP 的核心思想在于把游戏世界(包括代码)认为是数据的集合，在编写程序时要尽可能利用缓存同时避免 cache miss。
- Performance-Sensitive Programming
  基于 DOP 的思想来设计高性能的程序
  1. Reducing Order Dependency
     避免程序对于代码执行顺序的依赖。
  2. False Sharing in Cache Line
     不同的线程之间尽可能地相互独立。
  3. Branch Prediction
     具有相同分支的程序在一起执行，比如说可以通过对数据进行排序的方式来避免错误的分支预测。
     更通用的方法是按照业务逻辑对数据进行分组，每一组中只使用相同的函数进行处理。这样可以完全避免分支判断从而极大地提升程序性能。
- Performance-Sensitive Data Arrangements
  数据的组织方式对于程序性能也有巨大的影响。
  `AOS 和 SOA`
  ![alt text](image-105.png)
  当程序需要对数据进行访问时 AOS 往往会产生大量的 cache miss，因此在高性能编程中更推荐**使用 SOA 的组织方式。**
- ECS 架构(entity component system)
  ECS 架构中则使用了 entity 的概念将不同的组件组织起来。`entity 实际上只是一个 ID，用来指向一组 component`。而 ECS 架构中的 `component 则只包括各种类型的数据，不包含任何具体的业务逻辑。当需要执行具体的计算和逻辑时则需要调用 system 来修改 component 中的数据`。这样游戏中的数据可以集中到一起进行管理，从而极大地提升数据读写的效率。
  1. Unity DOTS 系统(data-oriented tech stacks)
  2. Unreal Mass System

![Performance](image-108.png)
执行速度

4.  Dynamic Global Illumination and Lumen
    `TODO`
5.  GPU-Driven Geometry Pipeline-Nanite
    `TODO`
6.  piccolo 代码讲解
    cmake 控制代码编译

- 代码结构
  **sourceTrail 工具看代码依赖**
  1. 小引擎使用全局变量 g_runtime global_context
     可以在任意地方方便访问各个系统,这种方式有什么优缺点?
     有没有其他方式能实现同样目的?
  2. GObject 上不同类型的 Component 调用顺序是不确定的
     比如一个 Player 可能先跑脚本再跑动画,一个先动画后脚本
     这会给游戏的运行带来什么问题?
     如何才能让 Component 有确定的顺序去执行?
  3. 这次视频我们只实现了 lua 脚本简单的执行。
     目前 lua 脚本每帧执行一次,而且还不能调用任何引擎的接口。
     你会怎么设计 luà 脚本的生命周期,以及与引擎各个系统的交互方式
- 反射系统
- 渲染系统
- PCG(Prodedual Content Generation，程序化内容生成) 迷宫生成算法
  PCG：使用随机来生成游戏中的内容；大量量产，例如生产地图。
  迷宫生成方法：并查集法
  ![alt text](image-109.png)

  n-1 条边联通 n 个结点，且无环
