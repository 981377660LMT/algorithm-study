1. 简单描述一下 Babel 的编译过程？
   1. 解析（Parse）：包括**词法分析**和**语法分析**。词法分析主要把字符流源代码（Char Stream）转换成令牌流（ Token Stream），语法分析主要是将令牌流转换成抽象语法树（Abstract Syntax Tree，AST）。
   2. 转换（Transform）：通过 Babel 的插件能力，将高版本语法的 AST 转换成支持低版本语法的 AST。当然在此过程中也可以对 AST 的 Node 节点进行优化操作，比如添加、更新以及移除节点等。
   3. 生成（Generate）：将 AST 转换成字符串形式的低版本代码，同时也能创建 Source Map 映射。

细节

- 解析： Babel 的解析过程（源码到 AST 的转换）可以使用 @babel/parser
- 转换: Babel 的转换过程（AST 到 AST 的转换）主要使用 @babel/traverse
  该库包可以通过**访问者模式**自动遍历并访问 AST 树的每一个 Node 节点信息，从而实现节点的替换、移除和添加操作，如下所示：
- Babel 的代码生成过程（AST 到目标代码的转换）主要使用 @babel/generator

```TS
import { parse } from '@babel/parser';
import traverse from '@babel/traverse';
import generate from '@babel/generator';

enum ParseSourceTypeEnum {
  Module = 'module',
  Script = 'script',
  Unambiguous = 'unambiguous',
}

enum ParsePluginEnum {
  Flow = 'flow',
  FlowComments = 'flowComments',
  TypeScript = 'typescript',
  Jsx = 'jsx',
  V8intrinsic = 'v8intrinsic',
}
const source = `let a: string = 1;`;

// 解析（Parser）阶段
const ast = parse(source, {
  // 严格模式下解析并且允许模块定义
  sourceType: ParseSourceTypeEnum.Module,
  // 支持解析 TypeScript 语法（注意，这里只是可以解析，并不是转换 TypeScript）
  plugins: [ParsePluginEnum.TypeScript],
});

// 转换（Transform) 阶段
traverse(ast, {
  // 访问词法规则
  VariableDeclaration(path) {
    path.node.kind = 'var';
  },

  // 访问词法规则
  TSTypeAnnotation(path) {
    // 移除 TypeScript 的声明类型
    path.remove();
  },
});

// 生成（Generate）阶段
const { code } = generate(ast);
// code:  var a = 1;
console.log('code: ', code);
```
