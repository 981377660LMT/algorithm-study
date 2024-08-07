一般，不推荐
https://wudashan.com/2017/05/03/Clean-Code-Read-Notes/
https://qianmo.blog.csdn.net/article/details/52079996?spm=1001.2014.3001.5502

- 把系统当故事来讲，而不是当做程序来写
- 整洁的代码简单直接，整洁的代码如优美的散文。整洁的代码从不隐藏设计者的意图，充满了干净利落的抽象和直截了当的控制语句
- 让代码比你来时更干净
  保持整洁的习惯，发现脏代码就要及时纠正。花时间保持代码代码整洁，这不但有关效率，还有关项目的生存
- 命名：
  - 准确、反映本质、直接
  - 避免造成误导:
    尽量提防长得太像的名称。想区分 XYZControllerForEfficientHandlingOfStrings 和 XYZControllerForEfficientStorageOfStrings，会花费我们太多的时间。
    以同样的方式拼写出同样的概念才是信息，拼写前后不一致就是误导。
  - 读得出来(讨论时不用解释)
  - 易于在 ide 中搜索
  - 每个概念对应唯一的一个词，一以贯之(例子：versionTree 的 mutation/operation 最后全变成了 step)
  - 使用解决方案领域专业术语
  - 添加有意义的语境
    只要短名称足够好，对含义的表达足够清除，就要比长名称更合适。添加有意义的语境甚好，别给名称添加不必要的语境
    **若没能提供放置的地方，还可以给名称添加前缀。**
- 函数参数中出现标识符参数是非常不推崇的做法。有标识符参数的函数，很有可能不止在做一件事，标示如果标识符为 true 将这样做，标识符为 false 将那样做。正确的做法应该将有标识符参数的函数一分为二，对标识符为 true 和 false 分别开一个函数来处理

- 优秀代码的格式准则
  https://qianmo.blog.csdn.net/article/details/52268975?spm=1001.2014.3001.5502

  - 像报纸一样一目了然(modules should be deep)
    优秀的源文件也要像报纸文章一样。
    名称应当`简单并且一目了然`，名称本身应该足够告诉`我们是否在正确的模块中`。
    源文件最顶部应该给出高层次概念和算法。细节应该往下渐次展开，直至找到源文件中最底层的函数和细节。
  - 恰如其分的注释
    带有少量注释的整洁而有力的代码，比带有大量注释的零碎而复杂的代码更加优秀。
    注释是为代码服务的，注释的存在大多数原因是为了代码更加易读，但`注释并不能美化糟糕的代码`
    注释存在的时间越久，就会离其所描述的代码的意义越远，越来越变得全然错误，`因为大多数程序员们不能坚持（或者因为忘了）去维护注释`
    当然，`教学性质的代码，多半是注释越详细越好`
  - 合适的单文件行数
    尽可能用几百行以内的单文件来构造出出色的系统
    《代码整洁之道》第五章中提到的 FitNess 系统，就是由大多数为 200 行、最长 500 行的单个文件来构造出总长约 5 万行的出色系统

  - 合理地运用空白行(@vue/reactivity 类型文件就做得很好)
    古诗中有`留白`，代码的书写中也要有适当的留白，也就是空白行。
    在每个命名空间、类、函数之间，都需用空白行隔开（应该大家在学编程之初，就早有遵守）。这条极其简单的规则极大地影响到了代码的视觉外观。每个空白行都是一条线索，标识出新的独立概念。
    其实，在往下读代码时，你会发现你的目光总停留于空白行之后的那一行。用空白行隔开每个命名空间、类、函数，代码的可读性会大大提升。
    **如果说空白行隔开了概念，那么靠近的代码行则暗示了他们之间的紧密联系**
  - 紧密相关的代码互相靠近

    ```java
    // 反例
      public class ReporterConfig
      {
          /**
          * The class name of the reporter listener
          */
          private String m_className;

          /**
          * The properties of the reporter listener
          */
          private List<Property> m_properties = new ArrayList<Property>();

          public void addProperty(Property property)
          {
              m_properties.add(property);
          }
      }

      // 正面示例
      public class ReporterConfig
      {
          private String m_className;
          private List<Property> m_properties = new ArrayList<Property>();

          public void addProperty(Property property)
          {
              m_properties.add(property);
          }
      }
    ```

    正面示例`一览无遗，一眼就能看这个是有两个变量和一个方法的类`
    反例，注释简直画蛇添足，`隔断了两个实体变量间的联系`，我们不得不移动头部和眼球，才能获得相同的理解度

  - 基于关联的代码分布
    - 大多数短函数，函数中的本地变量应当在函数的顶部出现
    - 若某个函数调用了另一个，就应该把它们放到一起，而且调用者应该尽量放到被调用者上面。这样，程序就有自然的顺序。若坚定地遵守这条约定，读者将能够确信函数声明总会在其调用后很快出现。
      概念相关的代码应该放到一起。相关性越强，则彼此之间的距离就该越短
      （和我的习惯不一样，public 和 private 分开)
    - `调用者应该尽量放到被调用者上面`(多用 function?)

- 整洁类的书写准则
  `代码质量与其整洁度成正比`，干净的代码，既在质量上可靠，也为后期维护、升级奠定了良好基础

  - 类中代码分布顺序·
    公有静态常量、私有静态常量、公有普通变量、私有普通变量、构造函数、公有函数、私有函数
    这样是符合自定向下的原则，`让程序读起来像一篇报纸文章`
  - 类或模块应有且只有一条加以修改的理由(单一职责、短小)
