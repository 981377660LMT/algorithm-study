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

   总结：

   - 通常实现的方式是游戏循环为每个演员每轮提供一小部分时间。每个实体都会有一个 Update()方法，该方法，执行一步然后返回；缺点是需要维护状态
   - 如果您的系统支持协程，由于`纤程维护自己的整个调用堆栈`，因此您甚至可以从其他函数调用中在它们之间进行切换

4. Amaranth, an Open Source Roguelike in C#
   Amaranth，一款用 C# 编写的开源 Roguelike 游戏

   - Because I’m crazy about decoupling, it’s actually split into three separate projects:
     因为我对解耦很着迷，所以它实际上分为三个独立的项目：
