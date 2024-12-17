1.  支持静态方法，在方法前使用 class 关键字来表示挂在类对象上的静态方法。

    **Smalltalk 和 Ruby 使用的 "元类 "是一种特别优雅的方法**
    `我们采用元类：每个类对象本身也是一个实例，因此它需要有自己的类，这个类就是元类。元类定义了可以在类对象上调用的方法，这些方法类似于其他语言中的静态方法。通过让 LoxClass 继承自 LoxInstance，类对象本身具备了实例的属性访问能力，这使得属性访问的逻辑可以统一处理实例属性和类属性（静态属性）。`

    Metaclasses are so cool, I almost wish the book itself discussed them
    properly, but there are only so many pages. `The idea is that a class object is itself an instance, which means it must have its own class -- a metaclass`. That metaclass defines the methods that are available on the
    class object -- what you'd think of as the "static" methods in a language like Java.

    Before we get to metaclasses, we need to push the new syntax through. In
    AstGenerator, add a new field to Stmt.Class:
    `在 Stmt.Class 中添加一个新的字段 classMethods 来存储静态方法。`

    ```java
    "Class      : Token name, List<Stmt.Function> methods, List<Stmt.Function> classMethods",
    ```

    When parsing a class, we separate out the class methods (prefixed with
    "class") into a separate list:
    在解析类时，识别并分离出静态方法和实例方法。

    ```java
    private Stmt classDeclaration() {
      Token name = consume(IDENTIFIER, "Expect class name.");

      List<Stmt.Function> methods = new ArrayList<>();
      List<Stmt.Function> classMethods = new ArrayList<>();
      consume(LEFT_BRACE, "Expect '{' before class body.");

      while (!check(RIGHT_BRACE) && !isAtEnd()) {
        boolean isClassMethod = match(CLASS);
        (isClassMethod ? classMethods : methods).add(function("method"));
      }

      consume(RIGHT_BRACE, "Expect '}' after class body.");

      return new Stmt.Class(name, methods, classMethods);
    }
    ```

    In the resolver, we need to make sure to resolve the class methods too:
    确保resolver中静态方法也被正确解析和处理。

    ```java
    for (Stmt.Function method : stmt.classMethods) {
      beginScope();
      scopes.peek().put("this", true);
      resolveFunction(method, FunctionType.METHOD);
      endScope();
    }
    ```

    They are resolved mostly like methods. They even have a "this" variable,
    which will be the class itself.

    Now we're ready for metaclasses. Change the declaration of LoxClass to:
    修改 LoxClass 的声明，`使其继承自 LoxInstance` 并实现 LoxCallable 接口。

    ```java
    class LoxClass extends LoxInstance implements LoxCallable {
      final String name;
      private final Map<String, LoxFunction> methods;

      LoxClass(LoxClass metaclass, String name,
            Map<String, LoxFunction> methods) {
        super(metaclass);
        this.name = name;
        this.methods = methods;
      }

      // ...
    }
    ```

    LoxClass now extends LoxInstance. Every class object is also itself an
    instance of a class, its metaclass. When we interpret a class declaration,
    we create two LoxClasses:

    `在解释类声明时，创建两个 LoxClass 对象：一个是元类，包含所有的静态方法，另一个是实际的类，包含实例方法。`
    元类的 metaclass 设置为 null，以避免无限递归。

    ```java
    @Override
    public Void visitClassStmt(Stmt.Class stmt) {
      environment.define(stmt.name.lexeme, null);
      Map<String, LoxFunction> classMethods = new HashMap<>();
      for (Stmt.Function method : stmt.classMethods) {
        LoxFunction function = new LoxFunction(method, environment, false);
        classMethods.put(method.name.lexeme, function);
      }

      LoxClass metaclass = new LoxClass(null,
          stmt.name.lexeme + " metaclass", classMethods);

      Map<String, LoxFunction> methods = new HashMap<>();
      for (Stmt.Function method : stmt.methods) {
        LoxFunction function = new LoxFunction(method, environment,
            method.name.lexeme.equals("init"));
        methods.put(method.name.lexeme, function);
      }

      LoxClass klass = new LoxClass(metaclass, stmt.name.lexeme, methods);
      environment.assign(stmt.name, klass);
      return null;
    }
    ```

    First, we create a metaclass containing all of the class methods. It has
    null for its metametaclass to stop the infinite regress. Then we create the
    main class like we did previously. The only difference is that we pass in
    the metaclass as its class.

    That's it. There are no other interpreter changes. Now that LoxClass is an
    instance of LoxInstance, the existing code for property gets now applies to
    class objects. On the last line of:

    ```lox
    class Math {
      class square(n) {
        return n * n;
      }
    }

    print Math.square(3); // Prints "9".
    ```

    The `.square` expression looks at the object on the left. It's a
    LoxInstance. We call `.get()` on that. That fails to find a field named
    "square" so it looks for a method on the object's class with that name. The
    object's class is the metaclass, and the method is found there. You can
    even put fields on classes now:
    当你通过类对象（如 Math）访问静态方法时，.get() 方法会先在类对象的字段中查找。如果找不到，就会在类对象的类（即元类）的方法中查找。

    ```lox
    Math.pi = 3.141592653;
    print Math.pi;
    ```

2.  扩展 Lox 以支持 getter 方法。这些方法在声明时没有参数列表。当访问具有该名称的属性时，将执行 getter 主体。

    ```lox
    class Circle {
      init(radius) {
        this.radius = radius;
      }

      area {
        return 3.141592653 * this.radius * this.radius;
      }
    }

    var circle = Circle(4);
    print circle.area; // Prints roughly "50.2655".
    ```

    The first implementation detail we have to figure out is how our AST
    distinguishes a getter declaration from the declaration of a method that
    takes no parameters. This is kind of cute, but we'll use a _null_
    parameter list to indicate the former and an _empty_ for the latter. So,
    when parsing a method (and only a method, there are no getter _functions_),
    we allow the parameter list to be omitted:

    `首先，需要修改 AST 以区分普通方法和 getter 方法。具体来说，当解析方法声明时，如果方法没有参数列表，则将其视为 getter 方法。`

    ```java
    private Stmt.Function function(String kind) {
      Token name = consume(IDENTIFIER, "Expect " + kind + " name.");

      List<Token> parameters = null;

      // 允许在 getter 方法中完全省略参数列表
      // 当 kind 为 "method" 时，允许省略参数列表，这样的方法被视为 getter 方法。
      //如果省略参数列表，parameters 将被设置为 null，否则为一个空列表或包含参数的列表。
      if (!kind.equals("method") || check(LEFT_PAREN)) {
        consume(LEFT_PAREN, "Expect '(' after " + kind + " name.");
        parameters = new ArrayList<>();
        if (!check(RIGHT_PAREN)) {
          do {
            if (parameters.size() >= 255) {
              error(peek(), "Can't have more than 255 parameters.");
            }

            parameters.add(consume(IDENTIFIER, "Expect parameter name."));
          } while (match(COMMA));
        }
        consume(RIGHT_PAREN, "Expect ')' after parameters.");
      }

      consume(LEFT_BRACE, "Expect '{' before " + kind + " body.");
      List<Stmt> body = block();
      return new Stmt.Function(name, parameters, body);
    }
    ```

    Now we need to make sure the rest of the interpreter doesn't choke on a
    null parameter list. We check for it when resolving:
    `需要确保解析器在解析函数和方法时能够正确处理 parameters 为 null 的情况。`

    ```java
    private void resolveFunction(Stmt.Function function, FunctionType type) {
      FunctionType enclosingFunction = currentFunction;
      currentFunction = type;

      beginScope();
      // 如果 function.params 为 null，表示这是一个 getter 方法，不需要定义参数。
      if (function.params != null) {
        for (Token param : function.params) {
          declare(param);
          define(param);
        }
      }
      resolve(function.body);
      endScope();
      currentFunction = enclosingFunction;
    }
    ```

    And when calling a LoxFunction:
    `在调用函数时，需要处理 parameters 为 null 的情况，以便正确执行 getter 方法。`

    ```java
    @Override
    public Object call(Interpreter interpreter, List<Object> arguments) {
      Environment environment = new Environment(closure);
      if (declaration.params != null) {
        for (int i = 0; i < declaration.params.size(); i++) {
          environment.define(declaration.params.get(i).lexeme,
              arguments.get(i));
        }
      }

      // ...
    }
    ```

    Now all that's left is to interpret getters. The only difference compared to methods is that the getter body is executed eagerly as soon as the property is accessed instead of waiting for a later call expression to invoke it.
    This isn't maybe the most elegant implementation, but it gets it done:
    `在解释器中，需要在访问属性时检测该属性是否是一个 getter 方法。如果是，则立即执行 getter 方法并返回其结果，而不是返回方法本身。`

    ```java
    @Override
    public Object visitGetExpr(Expr.Get expr) {
      Object object = evaluate(expr.object);
      if (object instanceof LoxInstance) {
        Object result = ((LoxInstance) object).get(expr.name);
        if (result instanceof LoxFunction &&
            ((LoxFunction) result).isGetter()) {
          result = ((LoxFunction) result).call(this, null);
        }

        return result;
      }

      throw new RuntimeError(expr.name,
          "Only instances have properties.");
    }
    ```

    After looking up the property, we see if the resulting object is a getter.
    If so, we invoke it right now and use the result of that. This relies on
    one little helper in LoxFunction:

    ```java
    public boolean isGetter() {
      return declaration.params == null;
    }
    ```

    And that's it.

3.  Python 和 JavaScript 允许你从对象自身的方法之外自由访问对象的字段。Ruby 和 Smalltalk 对实例状态进行了封装。只有类上的方法可以访问原始字段，而且由类决定暴露哪些状态。大多数静态类型语言都提供了像 private 和 public 这样的修饰符，用于控制类的哪些部分可以按成员从外部访问。这些方法之间有什么权衡，为什么一种语言会偏爱其中一种？

    Python and JavaScript allow you to freely access the fields on an object
    from outside of the methods on that object. Ruby and Smalltalk encapsulate
    instance state. Only methods on the class can access the raw fields, and it
    is up to the class to decide which state is exposed using getters and
    setters. Most statically typed languages offer access control modifiers
    like `private` and `public` to explicitly control on a per-member basis
    which parts of a class are externally accesible.

    What are the trade-offs between these approaches and why might a language
    might prefer one or the other?

    The decision to encapsulate at all or not is the classic
    trade-off between whether you want to make things easier for the class
    _consumer_ or the class _maintainer_. By making everything public and
    freely externally visible and modifier, a downstream user of a class has
    more freedom to pop the hood open and muck around in the class's internals.

    However, that access tends to increasing coupling between the class and its
    users. That increased coupling makes the class itself more brittle, similar
    to the "fragile base class problem". If users are directly accessing
    properties that the class author considered implementation details, they
    lose the freedom to tweak that implementation without breaking those users.
    The class can end up harder to change. That's more painful for the
    maintainer, but also has a knock-on effect to the consumer -- if the class
    evolves more slowly, they get fewer newer features for free from the
    upstream maintainer.

    On the other hand, free external access to class state is a simpler, easier
    user experience when the class maintainer and consumer are the same person.
    If you're banging out a small script, it's handy to be able to just push
    stuff around without having to go through a lot of ceremony and boilerplate.
    At small scales, most language features that build fences in the program are
    more annoying than they are useful.

    As the program scales up, though, those fences become increasingly important
    since no one person is able to hold the entire program in their head.
    Boundaries in the code let you make productive changes while only knowing a
    single region of the program.

    Assuming you do want some sort of access control over properties, the next
    question is how fine-grained. Java has four different access control levels.
    That's four concepts the user needs to understand. Every time you add a
    member to a class, you need to pick one of the four, and need to have the
    expertise and forethought to choose wisely. This adds to the cognitive load
    of the language and adds some mental friction when programming.

    However, at large scales, each of those access control levels (except maybe
    package private) has proven to be useful. Having a few options gives class
    maintainers precise control over what extension points the class user has
    access to. While the class author has to do the mental work to pick a
    modifier, the class _consumer_ gets to benefit from that. The modifier
    chosen for each member clearly communicates to the class user how the class
    is intended to be used. If you're subclassing a class and looking at a sea
    of methods, trying to figure out which one to override, the fact that one
    is protected while the others are all private or public makes your choice
    much easier -- it's a clear sign that that method is for the subclass's
    use.

    ### 不同语言的属性访问控制方式

    1. **自由访问属性的语言**：
       - **Python 和 JavaScript**：这些语言允许你从对象外部自由访问和修改对象的字段（属性）。例如，可以直接读取或赋值对象的属性，而无需通过类的方法。
    2. **封装实例状态的语言**：

       - **Ruby 和 Smalltalk**：这些语言对实例状态进行了封装。只有类的方法可以访问对象的原始字段，类通过定义 getter 和 setter 方法来决定哪些属性对外暴露。这种方式限制了外部对对象内部状态的直接访问。

    3. **静态类型语言的访问控制修饰符**：
       - **Java 等语言**：大多数静态类型语言提供了像 `private`、`public`、`protected` 等访问控制修饰符，用于明确控制类的成员在外部的可访问性。这些修饰符允许开发者在每个成员级别上细粒度地管理访问权限。

    ### 这些方法之间的权衡

    4. **对类的使用者（消费者）的影响**：

       - **自由访问**：
         - **优点**：使用者可以更灵活地操作对象，简化了代码编写，特别是在小规模项目或脚本中，减少了需要编写的样板代码（boilerplate）。
         - **缺点**：使用者可以直接访问和修改对象的内部状态，可能导致对类内部实现的依赖性增强。如果类的内部实现发生变化，可能会破坏依赖这些内部状态的代码，增加维护难度。
       - **封装和访问控制**：
         - **优点**：通过限制对内部状态的访问，降低了类与其使用者之间的耦合度，提高了类的可维护性和可扩展性。类的实现细节被隐藏，允许类的维护者在不影响使用者的情况下修改内部实现。
         - **缺点**：增加了代码编写的复杂性，使用者需要通过类的方法来访问和修改属性，可能导致更多的样板代码。

    5. **对类的维护者（开发者）的影响**：

       - **自由访问**：
         - **优点**：开发者可以快速开发和测试，减少了需要编写的代码量。
         - **缺点**：难以保证类的内部状态的一致性和正确性，增加了调试和维护的复杂性。
       - **封装和访问控制**：
         - **优点**：提供了更好的控制和保证类的内部状态的正确性，简化了维护工作。
         - **缺点**：需要更多的设计和规划，选择合适的访问修饰符增加了认知负担。

    6. **语言设计的权衡**：

       - **简洁性 vs. 控制性**：
         - **简洁性**：一些语言（如 Python 和 JavaScript）追求简洁和灵活，允许自由访问对象属性，适合快速开发和小型项目。
         - **控制性**：其他语言（如 Ruby、Smalltalk、Java）强调封装和访问控制，适合大型项目和需要高可维护性的系统。

    7. **访问控制的细粒度**：
       - **细粒度控制**：
         - **优点**：如 Java 提供多种访问控制级别（`private`、`protected`、`public`、包私有），允许开发者对每个成员进行精确的访问控制，提升代码的可读性和可维护性。
         - **缺点**：增加了语言的复杂性和学习曲线，开发者需要理解和正确使用不同的访问修饰符。

    ### 为什么语言会偏爱其中一种方式？

    8. **语言的设计哲学**：

       - **灵活性与简洁性**：如 Python 和 JavaScript 更注重开发者的灵活性和快速迭代，允许自由访问对象属性，适合动态和快速发展的项目需求。
       - **安全性与可维护性**：如 Ruby、Smalltalk 和 Java 更注重代码的安全性和可维护性，通过封装和访问控制减少潜在的错误和维护成本，适合大型和长期维护的项目。

    9. **目标用户群体**：

       - **快速开发者**：需要快速原型设计和迭代的开发者可能更倾向于使用允许自由访问属性的语言。
       - **企业级开发**：需要高度可靠性和可维护性的企业级应用更可能选择提供细粒度访问控制的语言。

    10. **项目规模和复杂性**：

        - **小型项目**：自由访问属性的方式更适合小型项目，减少了不必要的复杂性。
        - **大型项目**：封装和访问控制在大型项目中尤为重要，帮助管理复杂性和模块化代码结构。

    ### 总结

    不同语言在属性访问控制上的设计选择反映了它们在灵活性、简洁性与安全性、可维护性之间的不同权衡。自由访问属性的语言适合快速开发和小型项目，而通过封装和访问控制的语言则更适合需要高可靠性和可维护性的复杂系统。选择哪种方式取决于语言的设计目标、目标用户群体以及预期的应用场景。

    如果你在实际开发中需要在这两种方式之间做选择，应该根据项目的规模、复杂性以及团队的维护能力来决定采用哪种访问控制策略。
