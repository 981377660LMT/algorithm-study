1. 什么是字节码?采用字节码的好处是什么?
   在 Java 中，JVM 可以理解的代码就叫做字节码（即扩展名为 .class 的文件），它不面向任何特定的处理器，只面向虚拟机。Java 语言通过字节码的方式，在一定程度上解决了传统解释型语言执行效率低的问题，同时又保留了解释型语言可移植的特点。所以 Java 程序运行时比较高效，而且，由于字节码并不针对一种特定的机器，因此，Java 程序无须重新编译便可在多种不同操作系统的计算机上运行。
   Java 程序从源代码到运行一般有 3 步：
   .java => .class 字节码文件 => JVM 进行 JIT 转成机器码并保存
   这也解释了我们为什么经常会说 Java 是`编译与解释共存的语言`。

   Java 虚拟机（JVM）是运行 Java 字节码的虚拟机。JVM 有针对不同系统的特定实现（Windows，Linux，macOS），目的是使用相同的字节码，它们都会给出相同的结果。字节码和不同系统的 JVM 实现是 Java 语言“一次编译，随处可以运行”的关键所在。

2. java 泛型
   常用的通配符为： T，E，K，V，？

？ 表示不确定的 java 类型
T (type) 表示具体的一个 java 类型
K V (key value) 分别代表 java 键值中的 Key Value
E (element) 代表 Element (容器的值)

3. Java 中的几种基本数据类型是什么
   | 基本类型 | 位数 | 字节 | 默认值 |
   | :-------- | :--- | :--- | :------ |
   | `int` | 32 | 4 | 0 |
   | `short` | 16 | 2 | 0 |
   | `long` | 64 | 8 | 0L |
   | `byte` | 8 | 1 | 0 |
   | `char` | 16 | 2 | 'u0000' |
   | `float` | 32 | 4 | 0f |
   | `double` | 64 | 8 | 0d |
   | `boolean` | 1 | | false |
   6 种数字类型 ：`byte`、short、int、long、float、double
   1 种字符类型：char
   1 种布尔型：boolean。

   自动装箱与拆箱
   装箱：将基本类型用它们对应的引用类型包装起来；
   拆箱：将包装类型转换为基本数据类型；
   Java 基本类型的包装类的大部分都实现了常量池技术。Byte,Short,Integer,Long 这 4 种包装类默认创建了数值 [-128，127] 的相应类型的缓存数据，Character 创建了数值在[0,127]范围的缓存数据，Boolean 直接返回 True Or False。

   `记住：所有整型包装类对象之间值的比较，全部使用 equals 方法比较。`

4. 在一个静态方法内调用一个非静态成员为什么是非法的?
   这个需要结合 JVM 的相关知识，`静态方法是属于类的，在类加载的时候就会分配内存`，可以通过类名直接访问。`而非静态成员属于实例对象，只有在对象实例化之后才存在`，然后通过类的实例对象去访问。在类的非静态成员不存在的时候静态成员就已经存在了，此时调用在内存中还不存在的非静态成员，属于非法操作。

5. String StringBuffer 和 StringBuilder 的区别是什么? String 为什么是不可变的?
   String 类中使用 `final 关键字修饰字符数组来保存字符串，private final char value[]，所以 String 对象是不可变的。`
   StringBuffer 对方法加了同步锁或者对调用的方法加了同步锁，所以是线程安全的。StringBuilder 并没有对方法进行加同步锁，所以是非线程安全的。

   对于三者使用的总结：
   操作少量的数据: 适用 String
   单线程操作字符串缓冲区下操作大量数据: 适用 StringBuilder
   多线程操作字符串缓冲区下操作大量数据: 适用 StringBuffer

6. Object 类的常见方法总结

```JAVA
public final native Class<?> getClass()//native方法，用于返回当前运行时对象的Class对象，使用了final关键字修饰，故不允许子类重写。

public native int hashCode() //native方法，用于返回对象的哈希码，主要使用在哈希表中，比如JDK中的HashMap。
public boolean equals(Object obj)//用于比较2个对象的内存地址是否相等，String类对该方法进行了重写用户比较字符串的值是否相等。

protected native Object clone() throws CloneNotSupportedException//naitive方法，用于创建并返回当前对象的一份拷贝。一般情况下，对于任何对象 x，表达式 x.clone() != x 为true，x.clone().getClass() == x.getClass() 为true。Object本身没有实现Cloneable接口，所以不重写clone方法并且进行调用的话会发生CloneNotSupportedException异常。

public String toString()//返回类的名字@实例的哈希码的16进制的字符串。建议Object所有的子类都重写这个方法。

public final native void notify()//native方法，并且不能重写。唤醒一个在此对象监视器上等待的线程(监视器相当于就是锁的概念)。如果有多个线程在等待只会任意唤醒一个。

public final native void notifyAll()//native方法，并且不能重写。跟notify一样，唯一的区别就是会唤醒在此对象监视器上等待的所有线程，而不是一个线程。

public final native void wait(long timeout) throws InterruptedException//native方法，并且不能重写。暂停线程的执行。注意：sleep方法没有释放锁，而wait方法释放了锁 。timeout是等待时间。

public final void wait(long timeout, int nanos) throws InterruptedException//多了nanos参数，这个参数表示额外时间（以毫微秒为单位，范围是 0-999999）。 所以超时的时间还需要加上nanos毫秒。

public final void wait() throws InterruptedException//跟之前的2个wait方法一样，只不过该方法一直等待，没有超时时间这个概念

protected void finalize() throws Throwable { }//实例被垃圾回收器回收的时候触发的操作

```

7. Java 序列化中如果有些字段不想进行序列化，怎么办？
   对于不想进行序列化的变量，使用 transient 关键字修饰。
8. 既然有了字节流,为什么还要有字符流?
   问题本质想问：不管是文件读写还是网络发送接收，信息的最小存储单元都是字节，那为什么 I/O 流操作要分为字节流操作和字符流操作呢？

回答：字符流是由 Java 虚拟机将字节转换得到的，问题就出在这个过程还算是非常耗时，并且，如果我们`不知道编码类型就很容易出现乱码问题。所以， I/O 流就干脆提供了一个直接操作字符的接口，方便我们平时对字符进行流操作`。如果`音频文件、图片等媒体文件用字节流比较好，如果涉及到字符的话使用字符流比较好。`

例如:`golang 的 IO 库 封装了 OS 操作`

9. JDK 代理模式
   代理模式是一种比较好理解的设计模式。简单来说就是 `我们使用代理对象来代替对真实对象(real object)的访问，这样就可以在不修改原目标对象的前提下，提供额外的功能操作，扩展目标对象的功能。`

   `代理模式的主要作用是扩展目标对象的功能，比如说在目标对象的某个方法执行前后你可以增加一些自定义的操作。`
   代理模式有静态代理和动态代理两种实现方式，我们 先来看一下静态代理模式的实现。

   静态代理中，`我们对目标对象的每个方法的增强都是手动完成的（后面会具体演示代码）`，非常不灵活（比如接口一旦新增加方法，目标对象和代理对象都要进行修改）且麻烦(需要对每个目标类都单独写一个代理类)。 实际应用场景非常非常少，日常开发几乎看不到使用静态代理的场景。

   相比于静态代理来说，动态代理更加灵活。从 JVM 角度来说，动态代理是`在运行时动态生成类字节码，并加载到 JVM 中的。`Spring AOP、RPC 框架应该是两个不得不提的，它们的实现都依赖了动态代理。
   动态代理在我们日常开发中使用的相对较少，但是在框架中的几乎是必用的一门技术。学会了动态代理之后，对于我们理解和学习各种框架的原理也非常有帮助。
