/* eslint-disable no-lone-blocks */
/* eslint-disable no-console */

// TODO:
// 1. Dependens on interface, not class.
// 2. Service locator, not only components.

import { Scanner } from './Scanner'
import { Parser } from './Parser'
import { Interpreter } from './Interpreter'
import { IToken, TokenType } from './types'
import { RuntimeError } from './consts'
import { Resolver } from './Resolver'

export class TsLox {
  private _hadError = false
  private _hadRuntimeError = false

  /**
   * 在一个REPL会话中连续调用run()时重复使用同一个解释器.
   * 目前这一点没有什么区别，但以后当解释器需要存储全局变量时就会有区别。
   * 这些全局变量应该在整个REPL会话中持续存在。
   */
  private readonly _interpreter: Interpreter

  constructor() {
    this._interpreter = new Interpreter({ reportError: this._runtimeError.bind(this) })
  }

  run(source: string): void {
    const scanner = new Scanner(source, { reportError: this._error.bind(this) })
    const tokens = scanner.scanTokens()
    const parser = new Parser(tokens, { reportError: this._error.bind(this) })
    const statements = parser.parse()
    if (this._hadError || this._hadRuntimeError) return
    if (!statements) return
    const resolver = new Resolver(this._interpreter, { reportError: this._error.bind(this) })
    resolver.resolve(statements)
    if (this._hadError || this._hadRuntimeError) return
    this._interpreter.interpret(statements)
  }

  private _error(pos: number | IToken, message: string): void {
    if (typeof pos === 'number') {
      this._report(pos, '', message)
    } else if (pos.type === TokenType.EOF) {
      this._report(pos.line, 'at end', message)
    } else {
      this._report(pos.line, `at '${pos.lexeme}'`, message)
    }
  }

  private _report(line: number, where: string, message: string): void {
    console.log(`[line ${line}] Error ${where}: ${message}`)
    this._hadError = true
  }

  private _runtimeError(error: RuntimeError): void {
    console.log(`${error.message}\n[line ${error.token.line}]`)
    this._hadRuntimeError = true
  }
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  const lox = new TsLox()

  lox.run('print (1 + 2) / 3;')

  lox.run(
    `
      var a = 1;
      var b = 2;
      print a + b;
      `
  )

  lox.run(
    `
      var a = 1;
      print a = 2;
      print a;
      `
  )

  lox.run(
    `
      var a = "global a";
      var b = "global b";
      var c = "global c";
      {
        var a = "outer a";
        var b = "outer b";
        {
          var a = "inner a";
          print a;
          print b;
          print c;
        }
        print a;
        print b;
        print c;
      }
      print a;
      print b;
      print c;
      `
  )

  lox.run(
    `
    if (0>1) print "true";
    else print "false";
    `
  )

  lox.run(
    `
    print nil or 1;
    print "a" or 1;
    `
  )

  lox.run(
    `
    a = 1;
    while (a < 10) {
      print a;
      a = a + 1;
    }
    `
  )

  console.log('Fibonacci:')
  lox.run(
    `
    var a = 0;
    var temp;
    for (var b = 1; a < 10000; b = temp + b) {
      print a;
      temp = a;
      a = b;
    }
    `
  )

  lox.run(
    `
    fun sayHi(first, last) {
      print "Hi, " + first + " " + last + "!";
    }

    sayHi("Dear", "Reader");
    `
  )

  lox.run(
    `
    fun fib(n) {
      if (n <= 1) return n;
      return fib(n - 2) + fib(n - 1);
    }

    for (var i = 0; i < 20; i = i + 1) {
      print fib(i);
    }
    `
  )

  console.log('Closure:')
  lox.run(
    `
    fun makeCounter() {
      var i = 0;
      fun count() {
        i = i + 1;
        print i;
      }
    
      return count;
    }
    
    var counter = makeCounter();
    counter(); // "1".
    counter(); // "2".
    `
  )

  lox.run(
    `
    var a = "global";
    {
      fun showA() {
        print a;
      }
    
      showA();
      var a = "block";
      showA();
    }
    `
  )

  // lox.run(
  //   `
  //   fun bad() {
  //     var a = 1;
  //     var a = 2;
  //   }
  //   `
  // )

  // lox.run(
  //   `
  //   return 1;
  //   `
  // )

  lox.run(
    `
    class Breakfast {
      cook() {
        print "Cooking breakfast.";
      }
    }

    print Breakfast;
    `
  )

  lox.run(
    `
    class Bagel {
      cook() {
        print "Cooking bagel.";
      }
    }

    var bagel = Bagel();
    bagel.a = 1;
    print bagel.a;
    bagel.cook();

    class Egotist {
      speak() {
        print this;
      }
    }

    var method = Egotist().speak;
    method();
    `
  )

  console.log('init:')

  lox.run(
    `
    class Doughnut {
      init() {
        print this;
      }
    }

    Doughnut();
    `
  )

  lox.run(
    `
    class Doughnut {
      cook() {
        print "Fry until golden brown.";
      }
    }

    class BostonCream < Doughnut {}

    BostonCream().cook();
    `
  )

  lox.run(
    `
    class Doughnut {
      cook() {
        print "Fry until golden brown.";
      }
    }

    class BostonCream < Doughnut {
      cook() {
        super.cook();
        print "Pipe full of custard and coat with chocolate.";
      }
    }

    BostonCream().cook();
   `
  )
}
