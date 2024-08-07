https://journal.stuffwithstuff.com/archive/

1. C# Extension Methods
   无侵入式扩展方法

   - “Extension method” = “friendlier calling convention”

   - Reuse methods without inheritance
     重用方法而不继承
     `类似 rust 的 为某个 trait 实现方法`
     “if this class provides this capability, then it also has this capability”.
   - prefer static methods of a helper class to instance methods
     `更喜欢工具类的静态方法而不是实例方法`

2. Checking Flags in C# Enums
   I like C# enums and I also like using them as `bitfields`
3. Using an Iterator as a Game Loop
   使用迭代器作为游戏循环

   ```cpp
   void GameLoop()
   {
       while (mPlaying)
       {
           HandleUserInput();
           UpdateGameState();
           Render();
       }
   }
   ```

   - Separating out the UI 分离 UI
     没有显式调用 Render() ，而是让引擎在事情发生时引发事件（怪物移动、使用物品等）
     ```cpp
     void ProcessGame(UserInput input)
     {
         HandleUserInput(input);
         UpdateGameState();
     }
     ```
   - Enter iterators 输入迭代器

     ```cpp
     IEnumerable<bool> CreateProcessIterator()
     {
         while (true)
         {
             foreach (Entity entity in Entities)
             {
                 if (entity.NeedsUserInput)
                 {
                     yield return true;
                 }
                 else
                 {
                     entity.Move();
                 }
             }

             foreach (Item item in Items)
             {
                 item.Move();
             }
         }
     }
     ```

4. Naming Things in Code 在代码中命名事物
   When I’m designing software, I spend a lot of time thinking about names. For me, thinking about names is inseparable from the process of design. To name something is to define it.
   当我设计软件时，我会花很多时间思考名字。对我来说，对名字的思考，离不开设计的过程。命名某物就是定义它。

   If a class doesn’t represent something concrete, consider a metaphor.
   如果一个类不代表具体的东西，考虑一个隐喻。

   - Bad:
     - IncomingMessageQueue
     - CharacterArray
     - SpatialOrganizer
   - Good:
     - Mailbox
     - String
     - Map

   If you use a metaphor, use it consistently.
   如果你使用一个比喻，请始终如一地使用它。

   - Bad: Mailbox, DestinationID
   - Good: Mailbox, Address

   Name functions that return a Boolean (i.e. predicates) like questions.
   命名返回布尔值（即谓词）的函数，如问题。

   - Bad: list.Empty();
   - Good: list.IsEmpty();
     list.Contains(item);

   Name functions that return a value and don’t change state using nouns.
   **命名返回值,且不使用名词更改状态的函数(仅查询)。**

   - Bad: list.GetCount();
   - Good: list.Count();

   Don’t make the name redundant with an argument.
   **不要用参数使名称多余。**

   - Bad:
     list.AddItem(item);
     handler.ReceiveMessage(msg);
   - Good:
     list.Add(item);
     handler.Receive(msg);

   Don’t use “and” or “or” in a function name.
   **不要在函数名称中使用“and”或“or”。**
   如果在名称中使用连词，则该函数可能做得太多。
   将其分解成较小的部分并相应地命名。如果要确保这是一个原子操作，请考虑为整个操作创建一个名称，或者可能创建一个封装它的类。

   - Bad:
     mail.VerifyAddressAndSendStatus();
   - Good:
     mail.VerifyAddress();
     mail.SendStatus();

   **命名良好的代码更容易与其他程序员讨论，有助于传播代码知识；**
   一个带有良好名称部件的模块可以快速教会您它的作用。通过只阅读一小部分代码，您将快速构建整个系统的完整心理模型。如果它将某物称为“邮箱”，您将希望看到“邮件”和“地址”，而无需阅读它们的代码。
   另一方面，糟糕的名称会创建一堵不透明的代码墙，迫使您在脑海中煞费苦心地运行该程序，观察其行为，然后创建自己的私人命名法。“哦，DoCheck（）看起来正在遍历连接，看看它们是否都是实时的。我称之为 AreConnectionsLive（）“。这不仅速度慢，而且不可转移。

   **当我无法命名某物时，很有可能我试图命名的东西设计不佳。也许它试图一次做太多事情，或者错过了一个关键部分来完成它。**
   很难说我的设计是否好，但我发现我`做得不好的最可靠的指南之一是当名字不容易出现时。`
   当我现在设计时，我试着注意这一点。`一旦我对名字感到满意，我通常会对设计感到满意。`
