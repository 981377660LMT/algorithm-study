https://go7hic.github.io/A-Philosophy-of-Software-Design/#/

---

读完这本书，收获更多的不是具体 design 的方法，而是 mindset。

现在回去看本书的封面觉得很有意思，上方是杂乱的线，代表复杂的 implementation，而下方是整齐的线，代表 interface；这表示本书最重要的一个观点，module should be deep，一个 deep module 有简单的 interface，但是有复杂的 implementation，隐藏了很多 module 使用者不需要的信息。deep module 有助于减少 complexity，对于 module 的使用者来说，需要了解得更少，但是能获得更多的。

另一个收获就是 investment mindset。我工作中也常常会以不知道如何 design，觉得 design 太花时间，改动太大可能会破坏以往已经实现的功能为由，避免 design，停留在 feature-driven development 的舒适区。但长时间地注重实现功能的开发（tactical programming）会导致 complexity 的叠加，此时需要修改一个很小的部分，也常常需要花很多的时间。而 investment mindset 的想法是可以花 20%左右的时间去思考这些问题：为什么要实现这个功能；怎样实现这个 feature 才能让它好像一开始就在 design 中一样；如何用别的办法实现这个功能，和原方法相比有什么利弊；怎样才能让代码的读者更容易理解我的代码；有哪些代码是冗余的。可能我刚开始需要花更多的时间，但我相信长期的训练可以有益于自己对更复杂软件的实现。

investment mindset 是一种视角的转变，之前我 evaluate 我工作的完成可能是以单一的时间维度：如果我能在很快的时间让代码工作起来，就是好的；但是现在加入了 complexity 的维度：我要更多得考虑 design，考虑减少系统的复杂度，让之后 maintain 起来更方便，也要提升自己对 system abstraction 的理解。

investment mindset 也是一种长期主义的心态，我应该做一些长期有益的事情。做项目并不是很短暂的过程，所以需要写 unit test 方便之后项目的重构，写 comment & document 方便之后修改的时候参考。同时生活中也是一样，我没必要纠结于一时的得失因为它会过去；我不能逃避面对问题因为之后同样的问题可能会再出现；我不用焦虑于当下自己做不好一些事情，因为只要我 keep practice，总有一天自然而然就能做好。虽然好像是一些大道理，但是花一些时间想清楚，接下来只要去实践。

总之，这本书相比别的具体介绍 design pattern 之类的设计书而言，并不是非常实用，但它用一整本书讲了为什么我要 design，究竟有什么好处；我没有看过其他 design 类的书籍，但现在也有兴趣去探索一下相关的书籍。而且全书的逻辑非常清晰，例子很多，读起来也不吃力，是一本值得推荐的好书。
