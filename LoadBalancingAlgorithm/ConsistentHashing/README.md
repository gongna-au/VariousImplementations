## Consistent Hasing

Consistent Hashing 是一种分布式哈希方案，它的运行独立于分布式哈希表中的服务器或对象的数量。它为许多高流量的动态网站和 Web 应用程序提供支持。一种这样的系统类型，即为许多高流量动态网站和 Web 应用程序提供动力的分布式缓存，通常由分布式哈希的特殊情况组成。这些利用了一种称为一致性散列的算法。

什么是散列？

什么是“散列”？Merriam-Webster 将名词 hash 定义为“切碎的肉与土豆混合并变成褐色”，动词定义为“将（作为肉和土豆）切成小块”。因此，抛开烹饪细节不谈，哈希大致意味着“切碎和混合”——而这正是技术术语的来源。散列函数是一种将一段数据（通常描述某种对象，通常是任意大小）映射到另一段数据（通常是整数，称为散列码，或简称为散列）的函数。

例如，一些设计用于散列字符串的散列函数，输出范围为 0 .. 100，可以将字符串映射 Hello 到，例如，数字 57，Hasta la vista, baby 到数字 33，以及任何其他可能的字符串到该范围内的某个数字。由于可能的输入比输出多得多，因此任何给定的数字都会有许多不同的字符串映射到它，这种现象称为碰撞。

介绍哈希表（哈希图）？

想象一下，我们需要保留某个俱乐部所有成员的列表，同时能够搜索任何特定成员。我们可以通过将列表保存在数组（或链表）中来处理它，并且为了执行搜索，迭代元素直到找到所需的元素（例如，我们可能根据它们的名称进行搜索）。在最坏的情况下，这意味着检查所有成员（如果我们要搜索的成员是最后一个，或者根本不存在），或者平均检查其中的一半。在复杂性理论术语中，搜索将具有复杂性 O(n)，对于一个小列表来说它会相当快，但它会变得越来越慢，与成员数量成正比。

假设所有这些俱乐部成员都有一个 member ID，它恰好是反映他们加入俱乐部的顺序的序号。

假设 search byID 是可以接受的，我们可以将所有成员放在一个数组中，它们的索引与它们 ID 的 s 匹配（例如，一个成员 ID=10 将位于 10 数组中的索引处）。这将允许我们直接访问每个成员，根本不需要搜索。这将是非常有效的，事实上，尽可能高效，对应于可能的最低复杂度 O(1)，也称为常数时间。

但是，不可否认，我们的俱乐部会员 ID 场景有些做作。如果 IDs 是大的、非顺序的或随机数怎么办？或者，如果 ID 不接受搜索，我们需要按名称（或其他字段）搜索？保持我们的快速直接访问（或接近的东西），同时**能够处理任意数据集**和**限制较少的搜索条件**肯定会很有用。

**给我们想要被查找的数据集合中的每一个数据都分配一个 ID**这就是哈希函数来拯救的地方。可以使用合适的散列函数将任意数据块映射到整数，这将起到与我们俱乐部成员 ID 相似的作用，尽管有一些重要的区别（随机的，较大的）。

首先，一个好的散列函数通常具有很宽的输出范围（通常是一个 32 位或 64 位整数的整个范围），因此**为所有可能的索引构建一个数组要么不切实际，要么根本不可能，并且会大量浪费内存.** 为了克服这个问题，我们可以有一个合理大小的数组（比如说，只是我们期望存储的元素数量的两倍）并对哈希执行模运算以获得数组索引。因此，索引将是，数组的大小在 index = hash(object) mod N

**为每一个数据集合中的元素分配一个标志身份的 hashId 后，对哈希表（一个用来映射）的长度取模，得到用来查询 hashId 的 index,但是不同的元素在经过这两个操作，得到的 index 可能会重复，如果有冲突（比如说，如果冲突那么就 index 加 1），那么简单的直接索引访问将不起作用**，其次，对象哈希不会是唯一的（除非我们使用的是固定数据集和定制的完美哈希函数，但我们不会在这里讨论）。会有冲突（通过模运算进一步增加），因此简单的直接索引访问将不起作用。有几种方法可以解决这个问题，但一种典型的方法是将一个列表（通常称为存储桶）附加到每个数组索引以保存共享给定索引的所有对象。

> 简单的说，就是一个数据元素通过 Hash 函数得到标志自己身份的 hashID ,hashID 取 Mod 就知道了 hashID 在存储数组里面的索引，但是由于冲突的存在导致存储数组里面这个索引下存储的并不是该数据元素的 hashID

**存储桶解决冲突情况**但一种典型的方法是将一个列表（通常称为存储桶）附加到每个数组索引以保存共享给定索引的所有对象（如果 index 冲突，就继续把该元素放在桶里面）因此，我们有一个 size 为 N 的数组，每个条目都指向一个对象桶。要添加一个新对象，我们需要计算它的 hash modulo N，并检查结果索引处的存储桶，如果它不存在则添加该对象。要搜索对象，我们也这样做，只是查看存储桶以检查对象是否存在。这种结构称为哈希表，尽管*桶内的搜索是线性的*，但大小合适的哈希表每个桶应该有相当少的对象数量，从而导致几乎恒定的时间访问（平均复杂度为 O(N/k)，其中 k 是桶）。

对于复杂对象，散列函数通常不会应用于整个对象，而是应用于一个键。在我们的俱乐部成员示例中，每个对象可能包含多个字段（如姓名、年龄、地址、电子邮件、电话），但我们可以选择电子邮件作为键，只要把哈希函数应用于电子邮件。事实上，密钥不必是对象的一部分；存储键/值对很常见，其中键通常是相对较短的字符串，值可以是任意数据。在这种情况下，哈希表或哈希映射被用作字典，这就是一些高级语言实现对象或关联数组的方式。

横向扩展：分布式哈希？

在某些情况下，可能需要或希望将哈希表拆分为多个部分，由不同的服务器托管。这样做的主要动机之一是绕过使用单台计算机的内存限制，允许构建任意大的哈希表（给定足够的服务器）。在这种情况下，对象（及其密钥）分布在多个服务器之间，因此得名
一个典型的用例是内存缓存的实现，例如 Memcached。

这样的设置由一个缓存服务器池组成，这些服务器**托管许多键/值对**，用于**提供对**最初存储（或计算）在其他地方的**数据的快速访问**。例如，为了减少数据库服务器上的负载并同时提高性能，可以将应用程序设计为首先从缓存服务器获取数据，并且只有当数据不存在时——这种情况称为缓存未命中——求助于数据库，运行相关查询并使用适当的键缓存结果，以便下次需要时可以找到它。

现在，分配是如何进行的？使用什么标准来确定在哪些服务器中托管哪些密钥？

最简单的方法是取服务器数量的哈希模数。也就是说，池的大小在 server = hash(key) mod N。 N 为了存储或检索密钥，客户端首先计算哈希，应用 modulo N 操作。N 为了存储或检索密钥，客户端首先计算哈希，应用 modulo N 操作，并使用生成的索引联系适当的服务器（可能通过使用 IP 地址查找表）。请注意，用于密钥分发的散列函数在所有客户端中必须相同，但不必与缓存服务器内部使用的相同。
让我们看一个例子。假设我们有三台服务器 ,A 和 B，C 并且我们有一些带有哈希值的字符串键：
KEY HASH HASH mod 3
"john" 1633428562 2
"bill" 7594634739 0
"jane" 5000799124 1
"steve" 9787173343 0
"kate" 3421657995 2

客户想要检索 key 的值 john。它 hash modulo 3 是 2，所以它必须联系服务器 C。在那里找不到密钥，因此客户端从源中获取数据并添加它。池看起来像这样：
| A | B | C |
| :-- | --: | :----: |
| | | "john" |

接下来另一个客户端（或同一个客户端）想要检索 key 的值 bill。它 hash modulo 3 是 0，所以它必须联系服务器 A。在那里找不到密钥，因此客户端从源中获取数据并添加它。池现在看起来像这样：
| A | B | C |
| :----- | --: | :----: |
| "bill" | | "john" |

添加剩余密钥后，池如下所示：
| A | B | C |
| :------ | -----: | :----: |
| "bill" | "jane" | "john" |
| "steve" | | "kate" |

重新散列问题？

这种分配方案简单、直观且运行良好。也就是说，直到服务器数量发生变化。如果其中一台服务器崩溃或不可用会怎样？当然，需要重新分配密钥以解决丢失的服务器问题。如果将一台或多台新服务器添加到池中，这同样适用；需要重新分配密钥以包含新服务器。对于任何分配方案都是如此，但我们简单的模数分配的问题在于，当服务器数量发生变化时，大多数 hashes modulo N 都会发生变化，因此大多数密钥都需要移动到不同的服务器上。因此，即使删除或添加了单个服务器，所有密钥也可能需要重新散列到不同的服务器中。在我们之前的示例中，如果我们删除 server C，我们必须使用 hash modulo 2 而不是重新散列所有键 hash modulo 3，键的新位置将变为：
KEY HASH HASH mod 2
"john" 1633428562 0
"bill" 7594634739 1
"jane" 5000799124 0
"steve" 9787173343 1
"kate" 3421657995 1

> 系统使用一致性哈希的方式就可以减少迁移数据时的宕机时间

> 在 Consistent Hashing 中，当哈希表被调整大小时（例如，向系统中添加了一个新的缓存主机），只需要重新映射 k/n 个键，其中 k 是键的总数，n 是服务器的总数。

> 首先当我们在分布式系统里我们有很多的主机（服务器）可分离，我们的数据可分割存在不同的主机中，一致的哈希可以做到当一个主机被移除时，把被移除的主机中的资料给部分其他过去的主机。而当新增一个主机的时候，也只需要从部分主机转移资料，不用每个主机都转移。

```go


```

> 关键点：每个服务器都有许多不同类型的散列函数，使用不同的散列函数进行散列，那么一个被集中的后端 peer 他可能来自散列 120，也可能来自散列 250，还可能来自散列 350，服务器(后端节点)生活在哈希环的多个位置。如果增加一个服(后端节点)（意味着哈希环就要调整大小，相应的增加对应的 key ） (哈希环就要重新映射) 那么哈希环有多少个 key ,需要重新映射 key/n 个键