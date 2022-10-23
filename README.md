# VariousImplementations

### 各种算法的实现

##### 1.负载均衡算法

- (Random)随机

  > 从后端列表中任意挑选一个出来，这就是随机算法.正式的 Random 的代码要比上面的核心部分还复杂一点点。原因在于我们还需要达成另外两个设计目标：(1)线程安全(2)可嵌套

- (Round-Robin)轮循

  > 挨个选择一个 peer

- (WeightedRoundRobin)加权轮循

  > 是给每个 peer 增加一个权重，在平均化轮询的基础上加上这个权重来作为调节。最著名的加权轮询算法实现要论及 Nginx 和 LVS 了。Nginx 平滑加权轮询。

- (Smooth Weight Round Robin)LVS 平滑加权轮询

  > LVS (Linux Virtual Server)到底是什么东西，其实它是一种集群(Cluster)技术，采用 IP 负载均衡技术和基于内容请求分发技术。调度器具有很好的吞吐率，将请求均衡地转移到不同的服务器上执行，且调度器自动屏蔽掉服务器的故障，从而将一组服务器构成一个高性能的、高可用的虚拟服务器。整个服务器集群的结构对客户是透明的，而且无需修改客户端和服务器端的程序。

- Hashing (源地址散列)

  > 当 Browser session 方式被无状态模型取而代之，所以 hash 算法其实有点落伍了。无状态模型的本质是在 header 中带上 token，这个 token 将能够被展开为用户身份登录后的标识，从而等效于 browser session。说了个半天，无状态模型，例如 JWT 等，其原始的意图就是为了横向缩放服务器群。

- Consistent Hashing (一致性哈希)

  > 997 年 MIT 的 Karger 发表了所谓一致性 Hashing 的算法论文，其与传统的 hashCode 计算的关键性不同在于，一方面将 hashCode 约束为一个正整数（int32/uint32/int64 等等）一方面将正整数空间 [0, int.MaxValue] 视为一个可回绕的环形，即所谓的 Hash Ring，而待选择的 peers 均匀地分布在这个环形上，从而保证了每次选取时能够充分平滑地选取到每个 peer。至于选取时的下标值的计算方面是没有限定的，所以你可以在这个下标值的计算方案上加入可选策略。在负载均衡领域中的一致性 Hash 算法，又加入了 Replica 因子，它实际上就是在计算 Peer 的 hash 值时为 peer 的主机名增加一个索引号的后缀，索引号增量 replica 次，这样就得到了该 peer 的 replica 个副本。这就将原本 n 台 peers 的规模扩展为 n x Replica 的规模，有助于进一步提高选取时的平滑度。

- Least Connections (最少连接数)

  > 如果一个后端服务上的活动都连接数时全部后端中最少的，那么就选它了。

##### 2.限流算法

##### 3.LRU Cache

### Web Session Golang 实现

### System Design - Consistent Hashing（一致性哈希）

### System Design - Consistent Hashing（Caching）
