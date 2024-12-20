// https://blog.csdn.net/u011649400/article/details/128345213
// https://pages.cs.wisc.edu/~remzi/OSTEP/Chinese/42.pdf

// !Data Journaling(预写日志（write-ahead logging）)
// 我们回到最开始介绍的例子，即需要写入Db，B[v2]，I[v2]到磁盘的例子。前面我们的操作是让这些三个需要写入磁盘的操作，分别写入到磁盘中。那么现在，我们可以将他们打包在一起写入log (或者说journal)，如下图所示。这里TxB~TxE表示一个log entry，箭头表示这个journal区域可能包含多个log entry，目前我们只使用一个log entry。
// TxB~TxE也表示一系列操作的组成的一个事务(transaction)，TxB是这个事务的头，包含这些操作需要写入到磁盘的位置信息(例如I[v2]，B[v2]，Db各自在磁盘的位置，一共三个blocks的地址)，以及一个本次事务关联的transaction identifier (TID) 。TxE标识事务的结束位置，也包含和TxB相同的TID。我们写入log的操作称为journal write。
// 如果这个事务被顺利写入到磁盘，那么我们可以就根据TxB记录的磁盘的三个blocks的位置信息，将I[v2]，B[v2]，Db各自写入到磁盘的位置，我们将这个操作称为checkpointing。如果这三个blocks都成功写入了磁盘，那么我们说这个文件系统被成功checkpointed，用于记录log的journal区域(即上图)也可以被释放，给其他操作使用。
// 讨论: 对于前面的例子，我们需要写入TxB，I[v2]，B[v2]，Db，TxE总共5个blocks。如果让这5个blocks一个个按顺序写入，这显然会影响性能。如果这5个blocks异步写入依然会存在一致性问题，例如TxB，B[v2]，Db，TxE成功写入到磁盘后，系统就马上崩溃，I[v2]，还没有写入。因此我们可以使用折衷的办法，异步写入TxB，I[v2]，B[v2]，Db，等它们都成功后，再写入TxE。所以我们可以将journaling分解为3个操作:
//
// 更新文件系统的协议三阶段：
// 日志写入(Journal write)：将TxB以及对应的文件操作写入到事务中，然后让它们异步写入，以及等待它们全部完成。
// 日志提交(Journal Commit)：写入TxE，并等待完成。完成后，我们称为这个事务是committed。
// 加检查点(Checkpoint)：将事务中的数据，分别各自回写到各自的磁盘位置中。
//
// Recovery恢复: 我们利用上面的三个操作去描述文件系统是如何保证一致性的。
// 崩溃发生在Journal Commit完成前：那么文件系统可以丢掉之前写入的log。由于磁盘具体位置的bitmap，inodes，data blocks都没变，所以可以确保文件系统一致性。
// 崩溃发生在Journal Commit后，Checkpoint之前：那么文件系统在启动时候，可以扫描所有已经commited的log，然后针对每一个log记录操作进行replay，即recovery的过程中执行Checkpoint，将log的信息回写到磁盘对应的位置。这种操作也成为重做日志（redo logging）。。
// 崩溃发生在Checkpoint完成后：那无所谓，都已经成功回写到磁盘了，文件系统的bitmap、inodes、data blocks也能确保一致性。
// ————————————————
