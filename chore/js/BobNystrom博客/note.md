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
