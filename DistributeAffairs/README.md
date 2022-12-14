### 2PC（两阶段提交）

两阶段提交是一种保证分布式系统数据一致性的协议，现在很多数据库都是采用的两阶段提交协议来完成 分布式事务 的处理。

### 模拟如下场景

在两阶段提交中，主要涉及到两个角色，分别是协调者和参与者。
第一阶段：当要执行一个分布式事务的时候，事务发起者首先向协调者发起事务请求，然后协调者会给所有参与者发送 prepare 请求（其中包括事务内容）告诉参与者你们需要执行事务了，如果能执行我发的事务内容那么就先执行但不提交，执行后请给我回复。然后参与者收到 prepare 消息后，他们会开始执行事务（但不提交），并将 Undo 和 Redo 信息记入事务日志中，之后参与者就向协调者反馈是否准备好了。
第二阶段：第二阶段主要是协调者根据参与者反馈的情况来决定接下来是否可以进行事务的提交操作，即提交事务或者回滚事务。
比如这个时候 所有的参与者 都返回了准备好了的消息，这个时候就进行事务的提交，协调者此时会给所有的参与者发送 Commit 请求 ，当参与者收到 Commit 请求的时候会执行前面执行的事务的 提交操作 ，提交完毕之后将给协调者发送提交成功的响应。
而如果在第一阶段并不是所有参与者都返回了准备好了的消息，那么此时协调者将会给所有参与者发送 回滚事务的 rollback 请求，参与者收到之后将会 回滚它在第一阶段所做的事务处理 ，然后再将处理情况返回给协调者，最终协调者收到响应后便给事务发起者返回处理失败的结果。
