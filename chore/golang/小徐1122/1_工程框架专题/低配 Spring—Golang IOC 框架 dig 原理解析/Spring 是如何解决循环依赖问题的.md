Spring 通过 **三级缓存机制** 解决单例作用域下的循环依赖问题，但其处理存在一定限制。以下是具体原理和关键细节：

---

### **1. 三级缓存的作用**

Spring 使用三级缓存管理 Bean 的创建过程：

1. **一级缓存（`singletonObjects`）**  
   存储 **完全初始化好的单例 Bean**，可直接使用。
2. **二级缓存（`earlySingletonObjects`）**  
   存储 **提前暴露的 Bean 半成品**（已实例化但未填充属性）。
3. **三级缓存（`singletonFactories`）**  
   存储 **ObjectFactory**，用于生成 Bean 的早期引用（解决代理对象问题）。

---

### **2. 解决循环依赖的流程**

以 **Bean A → Bean B → Bean A** 的循环依赖为例：

1. **创建 Bean A**
   - 实例化 A（调用构造函数），此时 A 未设置属性，是一个半成品。
   - 将 A 的 ObjectFactory 存入 **三级缓存**（`singletonFactories`）。
2. **填充 A 的属性**
   - 发现 A 依赖 B，触发 **创建 Bean B**。
3. **创建 Bean B**
   - 实例化 B，将 B 的 ObjectFactory 存入三级缓存。
   - 填充 B 的属性时，发现依赖 A，**从三级缓存中获取 A 的 ObjectFactory**，生成 A 的早期引用（可能为代理对象）。
   - 将 A 的早期引用存入 **二级缓存**（`earlySingletonObjects`），并从三级缓存移除。
4. **完成 B 的初始化**
   - B 的属性填充完成后，B 被标记为完全初始化，存入 **一级缓存**。
5. **回到 A 的初始化**
   - A 通过一级缓存获取已初始化的 B，完成属性填充。
   - A 被标记为完全初始化，存入一级缓存，并从二级缓存移除。

---

### **3. 关键设计点**

- **提前暴露半成品**  
  允许在 Bean 未完全初始化前暴露引用，打破循环依赖的僵局。
- **ObjectFactory 的延迟处理**  
  三级缓存通过 ObjectFactory 延迟生成 Bean 的引用，确保 **AOP 代理对象的正确注入**（代理可能在初始化后生成）。
- **缓存升级机制**  
  Bean 从三级缓存 → 二级缓存 → 一级缓存逐级提升，确保依赖链中始终使用正确的引用。

---

### **4. 限制条件**

Spring **并非能解决所有循环依赖**，需满足以下条件：

1. **单例作用域（Singleton）**  
   原型（Prototype）作用域的 Bean 无法解决循环依赖，因每次请求都会创建新实例，无法通过缓存提前暴露。
2. **依赖注入方式为属性注入（Field）或 Setter 注入**  
   **构造器注入（Constructor）无法解决循环依赖**，因实例化前需完成所有依赖注入，导致无法提前暴露引用。
3. **无自定义 BeanPostProcessor 干扰**  
   若 BeanPostProcessor 在实例化阶段修改 Bean 的依赖关系，可能导致缓存机制失效。

---

### **5. 代码示例与异常**

- **构造器注入导致异常**

  ```java
  @Service
  public class A {
      private final B b;
      public A(B b) { this.b = b; }  // 构造器注入
  }

  @Service
  public class B {
      private final A a;
      public B(A a) { this.a = a; }  // 构造器注入
  }
  ```

  启动时会抛出 `BeanCurrentlyInCreationException`。

- **正确写法（属性注入）**

  ```java
  @Service
  public class A {
      @Autowired  // 属性注入
      private B b;
  }

  @Service
  public class B {
      @Autowired  // 属性注入
      private A a;
  }
  ```

  正常启动，Spring 通过三级缓存解决循环依赖。

---

### **6. 总结**

Spring 的循环依赖解决方案是工程上的权衡：

- **核心思想**：通过 **提前暴露半成品** + **缓存引用** 打破循环。
- **适用场景**：单例 Bean 的属性注入或 Setter 注入。
- **局限性**：无法处理构造器注入和原型作用域的循环依赖。

理解这一机制有助于在复杂依赖场景中规避设计问题，同时避免过度依赖框架特性导致代码难以维护。
