1. 倘若一个存在数据状态滞后的 follower 提供了服务，客户端就可能读到老数据怎么办？
   `ReadIndex机制`：若需要允许从Follower读取（如优化读性能），Raft可通过ReadIndex或Lease Read机制确保读取的数据至少是已提交的版本，避免脏读。

   - Leader向多数节点发送心跳确认自身仍是Leader。
   - Follower在响应读请求前，需确认其日志已同步到Leader的最新commitIndex。

2. 集群一主多从，纵览全局. 倘若 leader 出了问题，群龙无首，系统岂不是会分崩离析吗？
   Raft的选举机制能在毫秒级恢复新Leader，结合多数派存活原则，确保系统高可用。
