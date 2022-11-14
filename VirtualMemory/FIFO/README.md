# 页替换算法

操作系统为何要进行页面置换呢？这是由于操作系统给用户态的应用程序提供了一个虚拟的“大容量”内存空间，而实际的物理内存空间又没有那么大。所以操作系统就就“瞒着”应用程序，只把应用程序中“常用”的数据和代码放在物理内存中，而不常用的数据和代码放在了硬盘这样的存储介质上。如果应用程序访问的是“常用”的数据和代码，那么操作系统已经放置在内存中了，不会出现什么问题。但当应用程序访问它认为应该在内存中的的数据或代码时，如果这些数据或代码不在内存中，会产生页访问异常。
这时，操作系统必须能够应对这种页访问异常，即尽快**把应用程序当前需要的数据或代码放到内存中**来，重新执行应用程序产生异常的访存指令，如果把硬存的数据调入内存的时候，内存中没有足够的空间，那么这个时候就需要淘汰几块页面，然后才能把需要的数据装入到内存里面。如何判断内存中哪些是“常用”的页，哪些是“不常用”的页，把“常用”的页保持在内存中，在物理内存空闲空间不够的情况下，把“不常用”的页置换到硬盘上就是页替换算法着重考虑的问题。

好的页替换算法会导致页访问异常次数少，也就意味着访问硬盘的次数也少，从而使得应用程序执行的效率就高。

先进先出(First In First Out, FIFO)页替换算法：

关键词：**当程序是线性的访问地址空间的时候**
该算法总是淘汰最先进入内存的页，即选择在内存中驻留时间最久的页予以淘汰。只需把一个应用程序在执行过程中已调入内存的页按先后次序链接成一个队列，队列头指向内存中驻留时间最久的页，队列尾指向最近被调入内存的页。这样需要淘汰页时，从队列头很容易查找到需要淘汰的页。FIFO 算法只是在应用程序按线性顺序访问地址空间时效果才好，否则效率不高。因为那些常被访问的页，往往在内存中也停留得最久，结果它们因变“老”而不得不被置换出去。FIFO 算法的另一个缺点是，它有一种异常现象（Belady 现象），即在增加放置页的页帧的情况下，反而使页访问异常次数增多。

# 页面置换机制

如果要实现页面置换机制，只考虑页替换算法的设计与实现是远远不够的，还需考虑其他问题：

- 一个虚拟的页和硬盘的扇区之间的对应关系？
- 什么时候换入和换出？
- 如何设计数据结构以来支持页面替换算法？
- 如何完成页面的换入和换出操作？
- 哪些页可以换出？

### 1. 可以被换出的页面？

**存储内核代码和数据的页面不要被换出去，因为内核是很快的，不要让内核执行到某一个步骤的时候去等待把存储在硬盘里面的内核代码和数据调入内存，这个是很慢的，而内核是很快的。**

在操作系统的设计当中，一个基本的原则是：**并非所有的物理页都可以交换出去的**，只有**映射到用户空间**且**被用户程序直接访问**的页面才能被交换，这里面的原因是什么呢？操作系统是执行的关键代码，需要保证运行的高效性和实时性，如果在操作系统执行过程中，发生了缺页现象，则操作系统不得不等很长时间（硬盘的访问速度比内存的访问速度慢 2~3 个数量级），这将导致整个系统运行低效。而且，不难想象，处理缺页过程所用到的内核代码或者数据如果被换出，整个内核都面临崩溃的危险。

但在实验三实现的 代码中，我们只是实现了换入换出机制，还没有设计用户态执行的程序，所以在实验三中在仅仅通过执行 check_swap 函数为内核中分配一些页，模拟对内核对应的这些页面的访问，然后通过 do_pgfault 调用 swap_map_swappable 函数查询这些页面的访问情况，并调用相关函数，换出不常用的页到磁盘上。

### 2. 虚存中的页与硬盘上的扇区之间的映射关系？

如果一个页被置换到了硬盘上，那操作系统如何能简捷来表示这种情况呢？在数据结构的设计上，充分利用了页表中的 PTE 来表示这种情况：

- 当一个 PTE 用来描述一般意义上的物理页时，显然它应该维护权限关系和（与实际物理地址的映射关系）
- 但当它用来描述一个被置换出去的物理页时，它被用来维护与（swap 物理硬盘扇区的映射关系）
- 并且该 PTE 不应该由 MMU 将它解释成物理页映射(即没有 PTE_P 标记)，与此同时对应的权限则交由 mm_struct 来维护
- 当对位于该页的内存地址进行访问的时候，必然导致 page fault，希望操作系统能够根据 PTE 描述的 swap 项将相应的物理页重新建立起来，并根据虚存所描述的权限重新设置好 PTE 使得内存访问能够继续正常进行。

### 3. 执行换入换出的时机

- 需要一个数据结构来表示当前所有合法的虚拟地址空间的集合。
- 应用程序访问地址所在的页不在内存时，就会产生 page fault 异常
- 产生异常调用，（对应的 PTE 的高 24 位不为 0，而最低位为 0）将 swap 中的数据调入
- 即积极换出策略和消极换出策略。积极换出策略是指操作系统周期性地（或在系统不忙的时候）主动把某些认为“不常用”的页换出到硬盘上，从而确保系统中总有一定数量的空闲页存在，这样当需要空闲页时，基本上能够及时满足需求；
- 消极换出策略是指，只是当试图得到空闲页时，发现当前没有空闲的物理页可供分配，这时才开始查找“不常用”页面，并把一个或多个这样的页换出到硬盘上。

在实验三中的基本练习中，支持上述的第二种情况。对于第一种积极换出策略，即每隔 1 秒执行一次的实现积极的换出策略，可考虑在扩展练习中实现。对于第二种消极的换出策略，则是在当代码 调用 alloc_pages 函数获取空闲页时，此函数如果发现无法从物理内存页分配器获得空闲页，就会进一步调用 swap_out 函数换出某页，实现一种消极的换出策略。

### 4.数据结构的设计

内存中物理页使用情况的变量是基于数据结构 Page 的全局变量 pages 数组，pages 的每一项表示了计算机系统中一个物理页的使用情况。为了表示物理页可被换出或已被换出的情况，可对 Page 数据结构进行扩展：

### 5.总结

1. 一个函数用来判断：**是否有空闲页面？**
2. 一个具体的 （执行具体的策略下）的函数，用来执行换出操作。
3. 一个全局的结构体切片 存储着每一个物理页对应的结构体的具体的使用情况
4. 一个结构体用来表示一个位于磁盘上的 swap page
5. 一个页替换算法的框架 用来实现各类的算法 swap manager
6. 一个函数，记录页面访问情况的相关属性，swap manager 实现
7. 一个函数，挑选需要换出的页面，swap manager 实现
8. 一个 tick_event （计时器）可以用来实现积极的换页策略。
9. 一个函数用来生成物理页面

#### 补充：

完整的流程应该是：

1 .Special kernel symbols『表示发生缺页中断』
2 .『表示分配物理地址』
entry 0xc010002c (phys)
etext 0xc010962b (phys)
edata 0xc0122ac8 (phys)
3 .memory management (表示内存管理)
4 .check_alloc_page () 『表示检查是否内存已经满了』
5 .check_pgdir() 『 表示是否换出成功 』
6 .print() 『 打印当前的页表情况』
7 .page fault 『 表示发生缺页』
8 .check_pgfault()
9 .check_vmm()
10 .执行各种换页策略