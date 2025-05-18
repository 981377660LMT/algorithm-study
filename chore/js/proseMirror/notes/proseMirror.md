https://marijnhaverbeke.nl/blog/collaborative-editing.html
ProseMirror 协同编辑算法

- The Problem 问题
  协同编辑系统是多人可以同时处理文档的系统。
  系统确保多个用户的文档数据保持同步，单个用户所做的更改将发送给其他用户，并显示在他们的文档中。
  一种解决方案是加锁，允许用户锁定文档（或文档的一部分），从而完全防止数据竞争问题。但这迫使等待锁释放。我们不这样做。
  如果我们允许并发更新，我们会遇到这样的情况：用户 A 和用户 B 都做了一些事情，不知道另一个用户的操作，现在他们所做的这些事情必须协调。当操作编辑文档的不同部分时，这些操作可能根本没有交互，或者当它们尝试更改同一个单词时，交互非常频繁。

- Operational Transformation 算法

  很多研究都是关于真正的分布式系统，其中`一组节点之间交换消息，没有中央控制点`。解决该问题的经典方法称为 op，就是这样一种`分布式算法`。它定义了`一种描述更改的方法`，该更改具有两个属性：

  1.  You can transform changes relative to other changes.
  2.  No matter in which order concurrent changes are applied, you end up with the same document.
      > 结合律与交换律
      > An Operational Transformation (OT) based system applies local changes to the local document immediately, and broadcasts them to other users. Those users will transform and apply them when they get them. In order to know exactly which local changes a remote change should be transformed through, such a system also has to send along some representation of the state of the document at the time the change was made.
      > 基于操作转换 （OT） 的系统会立即将本地更改应用于本地文档，并将其广播给其他用户。这些用户在获得它们时将转换和应用它们。为了确切地知道远程更改应该通过哪些本地更改进行转换，这样的系统还必须发送一些文档在进行更改时状态的表示。

  - Centralization 中心化
    The design decisions that make the OT mechanism complex largely stem from the need to have it be truly distributed.
    使 OT 机制变得复杂的设计决策很大程度上源于使其真正分布的需要
    But you can save oh so much complexity by introducing a central point. I am, to be honest, extremely bewildered by Google's decision to use OT for their Google Docs—a centralized system.
    `但是你可以通过引入一个中心点来节省这么多的复杂性`。老实说，我对 Google 决定将 OT 用于他们的 Google Docs（一个集中式系统）感到非常困惑。
    ProseMirror's algorithm is centralized, in that it has a single node (that all users are connected to) making decisions about the order in which changes are applied. This makes it relatively easy to implement and to reason about.
    ProseMirror 的算法是集中式的，因为它有一个节点（所有用户都连接到该节点）来决定应用更改的顺序。
    And I don't actually believe that this property represents a huge barrier to actually running the algorithm in a distributed way. Instead of a central server calling the shots, you could use a consensus algorithm like Raft to pick an `arbiter`. (But note that I have not actually tried this.)
    而且我实际上`并不认为这代表了以分布式方式实际运行算法的巨大障碍`。您可以使用像 Raft 这样的共识算法来选择仲裁者，而不是由中央服务器发号施令。（但请注意，我实际上还没有尝试过。

  - The Algorithm 算法
    Like OT, ProseMirror uses a change-based vocabulary and transforms changes relative to each other. Unlike OT, it does not try to guarantee that applying changes in a different order will produce the same document.
    与 OT 一样，ProseMirror 使用基于变化的词汇表，并转换彼此相关的变化。与 OT 不同，`它不试图保证以不同的顺序应用更改将生成相同的文档`。
    By using a central server, it is possible—easy even—to have clients all apply changes in the same order. You can use a mechanism much like the one used in code versioning systems. When a client has made a change, they try to push that change to the server. If the change was based on the version of the document that the server considers current, it goes through. If not, the client must pull the changes that have been made by others in the meantime, and rebase their own changes on top of them, before retrying the push.
    `通过使用中央服务器，可以（甚至很容易）让所有客户端以相同的顺序应用更改。您可以使用一种与代码版本控制系统中使用的机制非常相似的机制。当客户端进行更改时，他们会尝试将该更改推送到服务器。如果更改基于服务器认为最新的文档版本，则该更改将通过。如果没有，客户端必须拉取其他人在此期间所做的更改，并在这些更改的基础上重新设置自己的更改，然后再重试推送。`
    Unlike in git, the history of the document is linear in this model, and a given version of the document can simply be denoted by an integer.
    与 git 不同，在这个模型中，`文档的历史记录是线性的，文档的给定版本可以简单地用整数表示`。
    Also unlike git, all clients are constantly pulling (or, in a push model, listening for) new changes to the document, and track the server's state as quickly as the network allows.
    与 git 不同的是，所有`客户端都在不断拉取（或者在推送模型中侦听）对文档的新更改，并在网络允许的情况下以最快的速度跟踪服务器的状态`。

    The only hard part is rebasing changes on top of others. This is very similar to the transforming that OT does. But it is done with the client's own changes, not remote changes.
    唯一困难的部分是`将更改重新建立在其他更改之上(rebase)`。`这与 OT 所做的转换非常相似。但它是通过客户自己的更改完成的，而不是远程更改`。

    Because applying changes in a different order might create a different document, rebasing isn't quite as easy as transforming all of our own changes through all of the remotely made changes.
    由于以不同的顺序应用更改可能会创建不同的文档，因此变基并不像通过所有远程更改转换我们自己的所有更改那么容易。

  - Position Mapping 位置映射
    Whereas OT transforms changes relative to other changes, ProseMirror transforms them using a derived data structure called a position map. Whenever you apply a change to a document, you get a new document and such a map, which you can use to convert positions in the old document to corresponding positions in the new document. The most obvious use case of such a map is adjusting the cursor position so that it stays in the same conceptual place—if a character was inserted before it, it should move forward along with the surrounding text.
    OT 相对于其他更改转换更改，而 `ProseMirror 使用称为位置映射的派生数据结构转换它们`。每当对文档应用更改时，都会获得一个新文档和这样的映射，您可以使用它来将旧文档中的位置转换为新文档中的相应位置。这种 map 最明显的用例是调整光标位置，使其保持在相同的概念位置 - 如果在它之前插入了一个字符，它应该与周围的文本一起向前移动。
    Transforming changes is done entirely in terms of mapping positions. This is nice—it means that we don't have to write change-type-specific transformation code. Each change has one to three positions associated with it, labeled from, to, and at. When transforming the change relative to a given other change, those positions get mapped through the other change's position map.
    转换更改完全根据映射位置来完成。这很好，这意味着我们不必编写特定于变更类型的转换代码。`每个更改都有一到三个与之关联的位置，标记为 from 、 to 和 at` 。在相对于给定的其他更改转换更改时，这些位置将通过其他更改的位置映射。
  - Rebasing Positions 变基
    To fix this, the system uses mapping `pipelines that are not just a series of maps, but also keep information about which of those maps are mirror images of each other`. When a position going through such a pipeline encounters a map that deletes the content around the position, the system scans ahead in the pipeline looking for a mirror images of that map. If such a map is found, we skip forward to it, and restore the position in the content that is inserted by this map, using the relative offset that the position had in the deleted content. A mirror image of a map that deletes content must insert content with the same shape.
    为了解决这个问题，系统使用的制图管线不仅仅是一系列贴图，而且还保留了有关哪些贴图是彼此的镜像的信息。当通过此类管道的位置遇到删除该位置周围内容的地图时，系统会在管道中向前扫描，以查找该地图的镜像。如果找到这样的地图，我们会向前跳到它，并使用该位置在已删除内容中的相对偏移量来恢复该地图插入的内容中的位置。删除内容的地图的镜像必须插入具有相同形状的内容。
  - Mapping Bias 映射偏差
    Whenever content gets inserted, a position at the exact insertion point can be meaningfully mapped to two different positions: before the inserted content, or after it. Sometimes the first is appropriate, sometimes the second. The system allows code that maps a position to choose what bias it prefers.
    每当插入内容时，可以有意义地将精确插入点的位置映射到两个不同的位置：插入内容之前或之后。有时第一个是合适的，有时是第二个。该系统允许映射位置的代码选择它喜欢的偏差。
  - Types of Changes 更改类型
    `An atomic change in ProseMirror is called a step`. Some things that look like single changes from a user interface perspective are actually decomposed into several steps. For example, if you select text and press enter, the editor will generate a delete step that removes the selected text, and then a split step that splits the current paragraph.
    `ProseMirror 中的原子更改称为Step`。从用户界面的角度来看，一些看起来像是单个更改的东西实际上被分解为几个步骤。例如，如果选择文本并按 Enter 键，编辑器将生成删除所选文本的删除步骤，然后生成拆分当前段落的拆分步骤。
    - addStyle/removeStyle
    - split
    - join
    - ancestor
    - replace
      replace is more complex than the other ones, and my initial impulse was to split it up into steps that remove and insert content. But because the position map created by a replace step needs to treat the step as atomic (positions have to be pushed out of all replaced content), I got better results with making it a single step.
      replace 比其他类型更复杂，我最初的冲动是将其拆分为删除和插入内容的步骤。但是，由于替换步骤创建的位置图需要**将该步骤视为原子**（必须将位置从所有替换内容中推出），因此我将其设置为单个步骤得到了更好的结果。
  - Intention 意图
    **解决协同冲突的一些约定**(prd?)

    I've tried to define the steps and the way in which they are rebased in so that rebasing yields unsurprising behavior. Most of the time, changes don't overlap, and thus don't really interact. But when they overlap, we must make sure that their combined effect remains sane.
    我试图定义步骤以及它们重新变基的方式，以便变基会产生不令人惊讶的行为。大多数时候，更改不会重叠，因此不会真正交互。但是，当它们重叠时，我们必须确保它们的综合效果保持理智。

    Sometimes a change must simply be dropped. When you type into a paragraph, but another user deleted that paragraph before your change goes through, the context in which your input made sense is gone, and inserting it in the place where the paragraph used to be would create a meaningless fragment.
    有时，必须简单地删除更改。当您键入一个段落，但另一个用户在您的更改完成之前删除了该段落时，`您的输入有意义的上下文将消失`，将其插入该段落曾经所在的位置将创建一个无意义的片段。

    If you tried to join two lists together, but somebody has added a paragraph between them, your change becomes impossible to execute (you can't join nodes that aren't adjacent), so it is dropped.
    如果您尝试将两个列表连接在一起，但有人在它们之间添加了一个段落，则您的更改将无法执行（您无法连接不相邻的节点），因此`它被删除`。

    In other cases, a change is modified but stays meaningful. If you made characters 5 to 10 strong, and another user inserted a character at position 7, you end up making characters 5 to 11 strong.
    在其他情况下，更改会被修改，但仍有意义。如果您将字符 5 到 10 设置为强，而另一个用户在位置 7 插入了一个字符，则最终将字符 5 到 11 强。

    And finally, some changes can overlap without interacting. If you make a word a link and another user makes it emphasized, both of your changes to that same word can happen in their original form.
    最后，一些更改可以在不交互的情况下重叠。如果`您将一个单词设置为链接，而另一个用户将其加重，则您对同一单词的两次更改都可以以原始形式发生`。

  - Offline Work 离线工作
    Silently reinterpreting or dropping changes is fine for real-time collaboration, where the feedback is more or less immediate—you see the paragraph that you were editing vanish, and thus know that someone deleted it, and your changes are gone.
    `对于实时协作来说，静默地重新解释或删除更改是可以的`，因为实时协作的反馈或多或少是即时的——你看到你正在编辑的段落消失了，因此知道有人删除了它，你的更改也消失了。

    For doing offline work (where you keep editing when not connected) or for a branching type of work flow, where you do a bunch of work and then merge it with whatever other people have done in the meantime, the model I described here is useless (as is OT). It might silently throw away a lot of work (if its context was deleted), or create a strange mishmash of text when two people edited the same sentence in different ways.
    对于离线工作（在没有连接时继续编辑）或分支类型的工作流，你做一堆工作，然后将其与其他人在此期间所做的任何事情合并，我在这里描述的模型是无用的（就像 OT 一样）。它可能会默默地扔掉很多工作（如果它的上下文被删除），或者当两个人以不同的方式编辑同一个句子时，会产生奇怪的文本大杂烩。

    In cases like this, I think a diff-based approach is more appropriate. `You probably can't do automatic merging—you have to identify conflicts had present them to the user to resolve. I.e. you'd do what git does`.
    `在这种情况下，我认为基于差异的方法更合适`。`您可能无法进行自动合并 - 您必须确定已将其呈现给用户解决的冲突。也就是说，你会做 git 所做的事情`，手动解决冲突？

  - Undo History 撤消历史记录
    How should the undo history work in a collaborative system? The widely accepted answer to that question is that `it definitely should not use a single, shared history`. If you undo, the last edit that `you` made should be undone, not the last edit in the document.
    在协作系统中，撤消历史记录应该如何运作？这个问题被广泛接受的答案是，它绝对不应该使用单一的、共享的历史。如果撤消，则应撤消`你`上次的编辑，而不是文档中的最后一次编辑。

    This means that the easy way to implement history, which is to simply roll back to a previous state, does not work. The state that is created by undoing your change, if other people's changes have come in after it, is a new one, not seen before.
    `这意味着实现历史记录的简单方法（即简单地回滚到以前的状态）不起作用`。如果其他人的更改是在更改之后出现的，那么通过撤消更改而创建的状态是一种新的状态，以前从未见过。

    To be able to implement this, I had to define changes (steps) in such a way that `they can be inverted, producing a new step that represents the change that cancels out the original step`.
    为了能够实现这一点，我必须以这样一种方式定义更改（步骤），以便它们可以反转，从而生成一个新步骤，`该步骤表示抵消原始步骤的更改`。

    `ProseMirror's undo history accumulates inverted steps`, and also keeps `track of all position maps between them and the current document version`. These are needed to be able to `map the inverted steps` to the current document version.
    `ProseMirror 的撤消历史记录会累积倒置步骤`，并跟踪它们与当前文档版本之间的所有位置图。需要这些步骤才能将反转步骤映射到当前文档版本。

    A downside of this is that if a user has made a change but is now idle while other people are editing the document, the position maps needed to move this user's change to the current document version pile up without bound. To address this, `the history periodically compacts itself`, mapping the inverted changes forward so that they start at the current document again. It can then discard the intermediate position maps.
    这样做的缺点是，如果用户进行了更改，但现在在其他人编辑文档时处于空闲状态，则将该用户的更改移动到当前文档版本所需的位置图将无限制地堆积起来。为了解决这个问题，`历史记录会定期压缩自己`，将反向更改向前映射，以便它们再次从当前文档开始。然后，它可以丢弃中间位置图。
