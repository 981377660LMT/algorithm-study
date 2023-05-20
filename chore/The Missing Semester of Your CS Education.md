计算机教育中缺失的一课
The Missing Semester of Your CS Education 中文版
https://missing-semester-cn.github.io/

## SHELL

文字接口：Shell
Bourne Again SHell, 简称 “bash” 。 这是被最广泛使用的一种 shell.
如果你要求 shell 执行某个指令，但是该指令并不是 shell 所了解的编程关键字，那么它会去咨询 环境变量 $PATH，它会列出当 shell 接到某条指令时，进行程序搜索的路径：

```shell
missing:~$ echo $PATH
/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

missing:~$ which echo
/bin/echo

missing:~$ /bin/echo $PATH
/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
```

通常，一个程序的输入输出流都是您的终端。也就是，您的键盘作为输入，显示器作为输出。 但是，我们也可以重定向这些流！
最简单的重定向是 < file 和 > file。这两个命令可以将程序的输入输出流分别重定向到文件。

## SHELL 脚本

Bash 中的字符串通过' 和 "分隔符来定义，但是它们的含义并不相同。以'定义的字符串为原义字符串，其中的变量不会被转义，而 `"定义的字符串会将变量值进行替换`。

## VIM

## 数据整理

## 命令行环境

某些情况下我们需要中断正在执行的任务，比如当一个命令需要执行很长时间才能完成时（假设我们在使用 find 搜索一个非常大的目录结构）。大多数情况下，我们可以使用 Ctrl-C 来停止命令的执行。但是它的工作原理是什么呢？为什么有的时候会无法结束进程？
当我们输入 Ctrl-C 时，shell 会发送一个 SIGINT 信号到进程。

### 配置文件（Dotfiles）

很多程序的配置都是通过纯文本格式的被称作点文件的配置文件来完成的（之所以称为点文件，是因为它们的文件名以 . 开头，例如 ~/.vimrc。也正因为此，它们默认是隐藏文件，ls 并不会显示它们）。

## 版本控制(Git)

尽管 Git 的接口有些丑陋，但是它的底层设计和思想却是非常优雅的。`丑陋的接口只能靠死记硬背，而优雅的底层设计则非常容易被人理解`。因此，我们将通过一种自底向上的方式向您介绍 Git。我们会从数据模型开始，最后再学习它的接口。一旦您搞懂了 Git 的数据模型，再学习其接口并理解这些接口是如何操作数据模型的就非常容易了。
可融合性可持久化数据结构
TrieNode

## 调试及性能分析

对于采样分析器来说，常见的显示 CPU 分析数据的形式是 火焰图，火焰图会在 Y 轴显示函数调用关系，并在 X 轴显示其耗时的比例。

## 元编程(这里其实主要讲的构建工具)

## 安全和密码学

熵(Entropy) 度量了不确定性并可以用来决定密码的强度。
熵的单位是 比特。对于一个均匀分布的随机离散变量，熵等于 log_2(所有可能的个数，即 n)。 扔一次硬币的熵是 1 比特。掷一次（六面）骰子的熵大约为 2.58 比特。
密码散列函数 (Cryptographic hash function) 可以将任意大小的数据映射为一个固定大小的输出。
一个散列函数的大概规范如下：

hash(value: array<byte>) -> vector<byte, N> (N 对于该函数固定)
`SHA-1`是 Git 中使用的一种散列函数， 它可以将任意大小的输入映射为一个 160 比特（可被 40 位十六进制数表示）的输出。
一个散列函数拥有以下特性：

- 确定性：对于不变的输入永远有相同的输出。
- 不可逆性：对于 hash(m) = h，难以通过已知的输出 h 来计算出原始输入 m。
- 目标碰撞抵抗性/弱无碰撞：对于一个给定输入 m_1，难以找到 m_2 != m_1 且 hash(m_1) = hash(m_2)。
- 碰撞抵抗性/强无碰撞：难以找到一组满足 hash(m_1) = hash(m_2)的输入 m_1, m_2（该性质严格强于目标碰撞抵抗性）。

### 对称加密

对称加密使用以下 3 个方法来实现这个功能：

```cpp
keygen() -> key (这是一个随机方法)
encrypt(plaintext: array<byte>, key) -> array<byte> (输出密文)
decrypt(ciphertext: array<byte>, key) -> array<byte> (输出明文)
```

加密方法 encrypt()输出的密文 ciphertext 很难在不知道 key 的情况下得出明文 plaintext。
解密方法 decrypt()有明显的正确性。因为功能要求给定密文及其密钥，解密方法必须输出明文：decrypt(encrypt(m, k), k) = m。
`AES` 是现在常用的一种对称加密系统。
对称加密的应用：加密不信任的云服务上存储的文件。对称加密和密钥生成函数配合起来，就可以使用密码加密文件：将密码输入密钥生成函数生成密钥 key = KDF(passphrase)，然后存储 encrypt(file, key)。

### 非对称加密

非对称加密的`“非对称”代表在其环境中，使用两个具有不同功能的密钥`： 一个是私钥(private key)，不向外公布；另一个是公钥(public key)，公布公钥不像公布对称加密的共享密钥那样可能影响加密体系的安全性。
非对称加密使用以下 5 个方法来实现加密/解密(encrypt/decrypt)，以及签名/验证(sign/verify)：

```cpp
keygen() -> (public key, private key)  (这是一个随机方法)
encrypt(plaintext: array<byte>, public key) -> array<byte>  (输出密文)
decrypt(ciphertext: array<byte>, private key) -> array<byte>  (输出明文)

sign(message: array<byte>, private key) -> array<byte>  (生成签名，类似`导师在表上签字`)
verify(message: array<byte>, signature: array<byte>, public key) -> bool  (验证签名是否是由和这个公钥相关的私钥生成的，类似验证`表上的签名是否由学生对应的导师签的`)
```

对称加密就好比一个防盗门：只要是有钥匙的人都可以开门或者锁门。
非对称加密好比一个可以拿下来的挂锁。你可以把打开状态的挂锁（公钥）给任何一个人并保留唯一的钥匙（私钥）。这样他们将给你的信息装进盒子里并用这个挂锁锁上以后，只有你可以用保留的钥匙开锁。

签名/验证方法具有和书面签名类似的特征。
在不知道 私钥 的情况下，不管需要签名的信息为何，很难计算出一个可以使 verify(message, signature, public key) 返回为真的签名。
对于使用私钥签名的信息，验证方法验证和私钥相对应的公钥时一定返回为真： verify(message, sign(message, private key), public key) = true。

非对称加密的应用

- PGP 电子邮件加密：用户可以将所使用的公钥在线发布，比如：PGP 密钥服务器或 Keybase。任何人都可以向他们发送加密的电子邮件。
- 聊天加密：像 Signal 和 Keybase 使用非对称密钥来建立私密聊天。
- 软件签名：Git 支持用户对提交(commit)和标签(tag)进行 GPG 签名。任何人都可以使用软件开发者公布的签名公钥验证下载的已签名软件。

## 大杂烩

- 修改键位映射
  将 Caps Lock 映射为 Ctrl 或者 Escape：Caps Lock 使用了键盘上一个非常方便的位置而它的功能却很少被用到，所以我们（讲师）非常推荐这个修改；
  将 PrtSc 映射为播放/暂停：大部分操作系统支持播放/暂停键；
- 守护进程
  以守护进程运行的程序名一般以 d 结尾，比如 SSH 服务端 sshd，用来监听传入的 SSH 连接请求并对用户进行鉴权。
- FUSE(用户空间文件系统)
  允许运行在用户空间上的程序实现文件系统调用，并将这些调用与内核接口联系起来。在实践中，这意味着用户可以在文件系统调用中实现任意功能。
- 备份
  推荐的做法是将数据备份到不同的地点存储。
  有效备份方案的几个核心特性是：`版本控制`，`删除重复数据`，以及`安全性`。
  对备份的数据实施版本控制保证了用户可以从任何记录过的历史版本中恢复数据。
  在备份中检测并删除重复数据，使其仅备份**增量变化**可以减少存储开销。
  在安全性方面，作为用户，你应该考虑别人需要有什么信息或者工具才可以访问或者完全删除你的数据及备份。最后一点，不要盲目信任备份方案。用户应该经常检查备份是否可以用来恢复数据。
- API（应用程序接口）
- 常见命令行标志参数及模式
  大部分工具支持 --help 或者类似的标志参数（flag）来显示它们的简略用法。
  **会造成不可撤回操作**的工具一般会提供“空运行”（dry run）标志参数，这样用户可以确认工具真实运行时会进行的操作。这些工具通常也会有“交互式”（interactive）标志参数，在执行每个不可撤回的操作前提示用户确认。
  --version 或者 -V 标志参数可以让工具显示它的版本信息（对于提交软件问题报告非常重要）。
  基本所有的工具支持使用 --verbose 或者 -v 标志参数来输出详细的运行信息。
  会造成破坏性结果的工具一般默认进行非递归的操作，但是支持使用“递归”（recursive）标志函数（通常是 -r）。

## 提问&回答

- 你会优先学习的工具有那些?

  多去使用键盘，少使用鼠标。
  学好编辑器。
  学习怎样去自动化或简化工作流程中的重复任务。
  学习像 Git 之类的版本控制工具并且知道如何与 GitHub 结合，以便在现代的软件项目中协同工作。

- 使用 Python VS Bash 脚本 VS 其他语言?

  对于大型或者更加复杂的脚本我们推荐使用更加成熟的脚本语言例如 Python 和 Ruby。

- 各种软件包和工具存储在哪里？引用过程是怎样的? /bin 或 /lib 是什么？

  根据你在命令行中运行的程序，这些包和工具会全部在 PATH 环境变量所列出的目录中查找到， 你可以使用 which 命令(或是 type 命令)来检查你的 shell 在哪里发现了特定的程序。
  /bin - 基本命令二进制文件
  /sbin - 基本的系统二进制文件，通常是 root 运行的
  /dev - 设备文件，通常是硬件设备接口文件
  /etc - 主机特定的系统配置文件
  /home - 系统用户的主目录
  /lib - 系统软件通用库
  /opt - 可选的应用软件
  /sys - 包含系统的信息和配置(第一堂课介绍的)
  /tmp - 临时文件( /var/tmp ) 通常重启时删除
  /usr/ - 只读的用户数据
  /usr/bin - 非必须的命令二进制文件
  /usr/sbin - 非必须的系统二进制文件，通常是由 root 运行的
  /usr/local/bin - 用户编译程序的二进制文件
  /var -变量文件 `像日志或缓存`,nginx 的日志文件在/var/log/nginx/下

- 我应该用 apt-get install 还是 pip install 去下载软件包呢?
- 用于提高代码性能，简单好用的性能分析工具有哪些?
- 你使用那些浏览器插件?

  全页屏幕捕获 GoFullPage - Full Page Screen Capture
  https://chrome.google.com/webstore/detail/gofullpage-full-page-scre/fdpohaocaechififmbbbbbknoalclacl?hl=en
  `Alt+Shift+P`

  密码集成管理器

- 有哪些有用的数据整理工具？
- Docker 和虚拟机有什么区别?
- 不同操作系统的优缺点是什么，我们如何选择（比如选择最适用于我们需求的- Linux 发行版）？
- 2FA 是什么，为什么我需要使用它?

  双因子验证（Two Factor Authentication 2FA）在密码之上为帐户增加了一层额外的保护。为了登录，你不仅需要知道密码，`还必须以某种方式“证明”可以访问某些硬件设备`。最简单的情形是可以通过接收手机的 SMS 来实现（尽管 SMS 2FA 存在 已知问题）。我们推荐使用 YubiKey 之类的 U2F 方案。

- 对于不同的 Web 浏览器有什么评价?
  2020 的浏览器现状是，`大部分的浏览器都与 Chrome 类似，因为它们都使用同样的引擎(Blink)`。Microsoft Edge 同样基于 Blink，而 Safari 则 基于 WebKit(与 Blink 类似的引擎)，这些浏览器仅仅是更糟糕的 Chrome 版本。不管是在性能还是可用性上，Chrome 都是一款很不错的浏览器。如果你想要替代品，我们推荐 Firefox。Firefox 与 Chrome 的在各方面不相上下，并且在隐私方面更加出色。 有一款目前还没有完成的叫 Flow 的浏览器，它实现了全新的渲染引擎，有望比现有引擎速度更快。