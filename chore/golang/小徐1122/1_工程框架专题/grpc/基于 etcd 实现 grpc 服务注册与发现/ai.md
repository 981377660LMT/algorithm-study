# 基于 etcd 实现 grpc 服务注册与发现

https://mp.weixin.qq.com/s/x-vC1gz7-x6ELjU-VYOTmA

我之前已经写过几篇关于gRPC服务端和客户端的源码分析，现在需要结合etcd来扩展这个主题。首先，我需要确认用户的需求是什么。他们可能是在实际项目中需要实现微服务的动态发现和负载均衡，所以想知道如何利用etcd来实现这一点。

首先，我需要回顾一下之前写的关于gRPC服务端和客户端的文章，确保这次的内容能够衔接上。用户提到使用etcd作为注册中心，并采用轮询负载均衡策略，所以重点应该放在服务注册、发现以及负载均衡的实现机制上。

接下来，我需要考虑文章的结构。用户提供的示例分为服务端和客户端两部分，分别涉及注册和发现。我需要详细解释这两个部分是如何工作的，包括etcd的租约机制、心跳维持、服务端注册的具体步骤，以及客户端如何通过etcd解析服务地址并应用负载均衡策略。

在服务端部分，用户展示了如何将服务节点注册到etcd，并定期续约。这里需要解释etcd的租约机制，为什么需要续约，以及如何通过etcd的API实现服务注册。可能还需要提到etcd的键值存储结构，如何以服务名为前缀存储节点信息。

客户端部分，用户使用了etcd的resolver来发现服务，并配置了轮询负载均衡。这里需要解释resolver的工作原理，如何监听etcd中的服务变化，并动态更新客户端的可用节点列表。同时，轮询策略如何在gRPC客户端中实现，如何确保每次请求选择不同的节点。

源码分析部分，需要深入gRPC客户端的resolver和balancer接口，以及etcd客户端库的相关实现。例如，etcd的resolver如何将服务地址转换为gRPC需要的格式，balancer如何根据策略选择节点。这部分可能需要结合代码示例，逐步讲解关键函数和流程。

另外，用户可能对etcd的watch机制感兴趣，需要解释etcd如何监控服务节点的变化，并在节点上下线时通知客户端更新列表。这是实现动态服务发现的核心，确保客户端总是获取最新的服务节点信息。

还需要考虑异常处理，比如etcd连接失败、节点失效等情况，gRPC客户端如何应对，是否有重试机制或故障转移策略。此外，负载均衡策略的扩展性，如何自定义其他策略如加权轮询或最少连接数。

最后，总结部分需要强调etcd在服务注册与发现中的优势，比如高可用、一致性，以及结合gRPC的负载均衡如何提升微服务架构的弹性和性能。可能还需要提到实际部署中的注意事项，比如etcd集群的配置、网络延迟的影响等。

---

### 基于 etcd 实现 gRPC 服务注册与发现详解

---

#### **一、核心架构**

1. **组件角色**：

   - **etcd**：分布式键值存储，提供数据一致性、租约机制和 Watch 监听功能。
   - **gRPC 服务端**：启动时注册自身信息到 etcd，定期续约维持存活状态。
   - **gRPC 客户端**：从 etcd 获取服务节点列表，动态监听变化，实现负载均衡。

2. **核心流程**：
   - **服务注册**：服务启动时写入 etcd，通过租约机制维持心跳。
   - **服务发现**：客户端从 etcd 获取节点列表，监听变更事件。
   - **负载均衡**：客户端基于策略（如轮询）选择节点发起请求。

---

#### **二、服务端实现**

1. **注册逻辑**：

   ```go
   func registerEndPointToEtcd(ctx context.Context, addr string) {
       etcdClient, _ := eclient.NewFromURL(MyEtcdURL)
       etcdManager, _ := endpoints.NewManager(etcdClient, MyService)
       lease, _ := etcdClient.Grant(ctx, 10) // 10s 租约
       _ = etcdManager.AddEndpoint(ctx,
           fmt.Sprintf("%s/%s", MyService, addr),
           endpoints.Endpoint{Addr: addr},
           eclient.WithLease(lease.ID))
       for {
           select {
           case <-time.After(5 * time.Second):
               etcdClient.KeepAliveOnce(ctx, lease.ID) // 续约
           case <-ctx.Done():
               return
           }
       }
   }
   ```

   - **键结构**：`服务名/节点地址`（如 `xiaoxu/demo/localhost:8080`）。
   - **租约机制**：节点需定期续约，超时后 etcd 自动删除节点键，实现自动下线。

2. **etcd 交互**：
   - **AddEndpoint**：将节点信息写入 etcd，绑定租约。
   - **KeepAlive**：定时续约，确保节点存活状态。
   - **优雅退出**：服务终止时自动释放租约，触发节点删除。

---

#### **三、客户端实现**

1. **服务发现配置**：

   ```go
   func main() {
       etcdClient, _ := eclient.NewFromURL("http://localhost:2379")
       etcdResolverBuilder, _ := eresolver.NewBuilder(etcdClient)
       etcdTarget := fmt.Sprintf("etcd:///%s", MyService)
       conn, _ := grpc.Dial(
           etcdTarget,
           grpc.WithResolvers(etcdResolverBuilder),
           grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy": "round_robin"}`),
       )
       client := proto.NewHelloServiceClient(conn)
   }
   ```

   - **Resolver**：通过 `etcd:///` 前缀声明使用 etcd 服务发现。
   - **负载均衡策略**：配置为轮询（`round_robin`）。

2. **etcd Resolver 工作原理**：

   - **Watch 监听**：通过 etcd 的 Watch API 监听服务节点变化。
   - **地址转换**：将 etcd 中的键值转换为 gRPC 的 `resolver.Address` 列表。
   - **动态更新**：节点变化时，触发 gRPC 客户端更新连接池。

   ```go
   func (r *resolver) watch() {
       for {
           select {
           case ups := <-r.wch:
               addrs := convertToGRPCAddress(ups) // 转换地址
               r.cc.UpdateState(gresolver.State{Addresses: addrs}) // 更新客户端
           }
       }
   }
   ```

---

#### **四、负载均衡机制**

1. **Round-Robin 策略**：

   - **Picker 实现**：维护节点列表和当前索引，每次请求轮询选择。

   ```go
   type rrPicker struct {
       subConns []balancer.SubConn
       next     uint32 // 原子计数器
   }
   func (p *rrPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
       index := atomic.AddUint32(&p.next, 1) % uint32(len(p.subConns))
       return balancer.PickResult{SubConn: p.subConns[index]}, nil
   }
   ```

   - **均衡性**：确保请求均匀分布到所有可用节点。

2. **节点状态管理**：
   - **健康检查**：gRPC 自动维护连接状态（如就绪、连接中、故障）。
   - **故障剔除**：节点不可达时，负载均衡器自动跳过。

---

#### **五、关键源码解析**

1. **etcd Resolver 构建**：

   ```go
   func NewBuilder(client *clientv3.Client) (gresolver.Builder, error) {
       return builder{c: client}, nil
   }
   func (b builder) Build(target gresolver.Target, cc gresolver.ClientConn, opts gresolver.BuildOptions) (gresolver.Resolver, error) {
       r := &resolver{c: b.c, target: target.Endpoint, cc: cc}
       r.wch, _ = endpoints.NewManager(r.c, r.target).NewWatchChannel(r.ctx)
       go r.watch() // 启动监听协程
       return r, nil
   }
   ```

   - **WatchChannel**：通过 etcd 的 Watch 接口监听服务前缀键。

2. **负载均衡器更新**：
   ```go
   func (b *baseBalancer) UpdateClientConnState(s balancer.ClientConnState) error {
       for _, addr := range s.ResolverState.Addresses {
           if _, ok := b.subConns.Get(addr); !ok {
               sc, _ := b.cc.NewSubConn([]resolver.Address{addr}, balancer.NewSubConnOptions{})
               sc.Connect() // 创建新连接
           }
       }
       b.regeneratePicker() // 重建 Picker
       b.cc.UpdateState(balancer.State{Picker: b.picker})
   }
   ```
   - **动态增删**：根据 etcd 返回的地址列表动态调整连接池。

---

#### **六、生产实践建议**

1. **etcd 集群部署**：至少 3 节点保障高可用，避免单点故障。
2. **租约时间优化**：根据网络状况调整 TTL，平衡心跳开销与故障检测速度。
3. **负载均衡策略**：
   - **权重轮询**：根据节点配置分配不同权重。
   - **地域优先**：优先选择同机房节点。
4. **异常处理**：
   - **重试机制**：配置 gRPC 的 `WithDefaultServiceConfig` 实现请求重试。
   - **熔断降级**：集成 Hystrix 或 Resilience4j 防止雪崩。

---

#### **七、总结**

- **服务注册**：通过 etcd 租约实现节点自动注册与心跳维持。
- **服务发现**：利用 etcd Watch 机制动态感知节点变化。
- **负载均衡**：gRPC 客户端内置策略（如轮询）实现请求分发。
- **扩展性**：可通过自定义 Resolver 和 Balancer 支持复杂场景。

通过结合 etcd 的强一致性与 gRPC 的高性能通信，可构建出弹性、可扩展的微服务架构。
