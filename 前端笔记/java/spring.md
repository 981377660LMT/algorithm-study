语言的流行通常需要一个杀手级的应用，Spring 就是 Java 生态的一个杀手级的应用框架。

1. 我们一般说 Spring 框架指的都是 Spring Framework，它是很多模块的集合，使用这些模块可以很方便地协助我们进行开发。
2. Spring 提供的核心功能主要是 IoC/DI 和 AOP。学习 Spring ，一定要把 IoC 和 AOP 的核心思想搞懂！
   谈谈自己对于 Spring IoC 的了解

   IoC（Inverse of Control:控制反转） 是一种设计思想，而不是一个具体的技术实现。IoC 的思想就是将原本在程序中手动创建对象的控制权，交由 Spring 框架来管理。不过， IoC 并非 Spirng 特有，在其他语言中也有应用。
   在实际项目中一个 Service 类可能依赖了很多其他的类，假如我们需要实例化这个 Service，你可能要每次都要搞清这个 Service 所有底层类的构造函数，这可能会把人逼疯。`如果利用 IoC 的话，你只需要配置好，然后在需要的地方引用就行了`，这大大增加了项目的可维护性且降低了开发难度。
   在 Spring 中， IoC 容器是 Spring 用来实现 IoC 的载体， `IoC 容器实际上就是个 Map（key，value），Map 中存放的是各种对象。`
   `Spring 时代我们一般通过 XML 文件来配置 Bean，后来开发人员觉得 XML 文件来配置不太好，于是 SpringBoot 注解配置就慢慢开始流行起来。`

   谈谈自己对于 AOP 的了解
   AOP(Aspect-Oriented Programming:面向切面编程)能够将那些与业务无关，却为业务模块所共同调用的逻辑或责任（例如事务处理、日志管理、权限控制等）封装起来，便于减少系统的重复代码，降低模块间的耦合度，并有利于未来的可拓展性和可维护性。
   Spring AOP 属于`运行时增强`，而 AspectJ `是编译时增强`。 Spring AOP `基于代理`(Proxying)，而 AspectJ 基于`字节码操作`(Bytecode Manipulation)。

   Spring bean
   什么是 bean？
   简单来说，bean 代指的就是那些被 IoC 容器所管理的对象。
   `我们需要告诉 IoC 容器帮助我们管理哪些对象`，这个是通过配置元数据来定义的。配置元数据可以是 XML 文件、注解或者 Java 配置类。

   ```XML
   <!-- Constructor-arg with 'value' attribute -->
   <bean id="..." class="...">
      <constructor-arg value="..."/>
   </bean>
   ```

3. Spring 框架中用到了哪些设计模式？(同样适用于 nestjs)
   - 工厂设计模式 : Spring 使用工厂模式通过 BeanFactory、ApplicationContext 创建 bean 对象。
   - 代理设计模式 : Spring AOP 功能的实现。
   - 单例设计模式 : Spring 中的 Bean 默认都是单例的。
   - 模板方法模式 : Spring 中 jdbcTemplate、hibernateTemplate 等以 Template 结尾的对数据库操作的类，它们就使用到了模板模式。
   - 外观模式 : 我们的项目需要连接多个数据库，而且不同的客户在每次访问中根据需要会去访问不同的数据库。这种模式让我们可以根据客户的需求能够动态切换不同的数据源。
   - 观察者模式: Spring 事件驱动模型就是观察者模式很经典的一个应用。
   - 适配器模式 : Spring AOP 的增强或通知(Advice)使用到了适配器模式、spring MVC 中也是用到了适配器模式适配 Controller。
   - 装饰者模式(AOP)
