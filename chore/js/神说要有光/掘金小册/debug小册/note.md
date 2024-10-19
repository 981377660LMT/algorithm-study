https://github.com/981377660LMT/articles/blob/af32bed62ae47d7c245d2089ce60f09f715b6b66/0706/debug%E5%B0%8F%E5%86%8C.md?plain=1#L4

## 总结

> 调试是接近原理最高效的方式

0. 定义
   调试就是通过某种信道（比如 WebSocket）把运行时信息传递给开发工具，做 UI 的展示和交互，辅助开发者排查问题、了解代码运行状态等
1. 调试工具四要素
   `frontend、backend、调试协议、信道`
   `Chrome DevTools 被设计成了和 Chrome 分离的架构，两者之间通过 WebSocket 通信，设计了专门的通信协议： Chrome DevTools Protocol。`
   Chrome DevTools 里展示的所有内容都是从 backend 那里拿到的，他只是一个展示和交互的 UI 而已。
   这个 UI 是可以换的，比如我们可以用 VSCode Debugger 对接 CDP 调试协议来调试网页。
   VSCode 不是 JS 专用编辑器，所以多了一个适配器层(Debugger Adapter)。

## vscode debugger 两个优势：

1. 回溯
2. 边调试边写代码是我推荐的写代码方式

## vscode debugger 配置笔记

1. userDataDir

`user data dir 是保存用户数据的地方，比如浏览历史、cookie 等，一个数据目录只能跑一个 chrome，所以默认会创建临时用户数据目录，想用默认的目录可以把这个配置设为 false`
![代价是，只能开一个浏览器实例](image.png)
用户数据目录有个特点，就是只能被一个 Chrome 实例所访问，如果你之前启动了 Chrome 用了这个默认的 user data dir，那就不能再启动一个 Chrome 实例用它了。
默认是 true，代表创建一个临时目录来保存用户数据。
你也可以指定一个自定义的路径，这样用户数据就会保存在那个目录下：

2. sourceMapPathOverrides
   ![alt text](image-1.png)

   `编译后代码路径 -- sourcemap -> 源码路径 -- sourceMapPathOverrides -> 本地文件路径`
   把调试的文件 sourcemap 到的路径映射到本地的文件，这样调试的代码就不再只读了：

   ```json
   // 默认把 meteor、webpack 开头的 path 映射到了本地的目录下
   // 其中 ?:* 代表匹配任意字符，但不映射，而 * 是用于匹配字符并映射的。
   "sourceMapPathOverrides": {
     "meteor://💻app/*": "${workspaceFolder}/*",
     "webpack:///./~/*": "${workspaceFolder}/node_modules/*",
     "webpack://?:*/*": "${workspaceFolder}/*"
   }
   ```

3. console
   internalConsole 就是内置的 debug console 面板，默认是这个。
   **internalTerminal 是内置的 terminal 面板，切换成这个就是彩色了** <- 推荐
   externalTerminal 会打开系统的 terminal 来展示日志信息：
4. env/envFile
   这样在调试的 node 程序里就可以取到这些环境变量
5. cwd
   指定工作目录，使用场景：运行对应目录下的 package.json 里的脚本
6. resolveSourceMapLocations
   **默认值是排除掉了 node_modules 目录的**，也就是不会查找 node_modules 下的 sourcemap。

## sourcemap 相关

开发时会使用 sourcemap 来调试，但是生产可不会，但是线上报错的时候确实也需要定位到源码，这种情况一般都是单独上传 sourcemap 到错误收集平台。
比如 sentry 就提供了一个  @sentry/webpack-plugin  支持在**打包完成后把 sourcemap 自动上传到 sentry 后台，然后把本地 sourcemap 删掉。还提供了  @sentry/cli  让用户可以手动上传。**

如何生成 sourcemap?
`用 source-map 库生成 sourcemap，`

## 调试 Vue 项目

1. 调试 @vue/cli 创建的 webpack 项目
   vue cli 创建的项目，默认情况下打断点不生效，这是因为文件路径后带了 `?hash`，这是默认的 eval-cheap-module-source-map 的 devtool 配置导致的，去掉 eval，改为 source-map 即可。
2. 调试 create vue 创建的 vite 项目
   `如果 sourcemap 到的文件路径不是本地路径，那就映射不到本地文件，导致断点打不上，这时候可以配置下 sourceMapPathOverrides。如果映射之后路径开头多了几层目录，那就要配置下 webRoot。`
   如果想点击调用栈直接在 workspace 打开对应的文件，这需要把 demo 项目和 vue3 源码项目放到一个 workspace 下，再次调试就可以了。

## VSCode Chrome Debugger 断点映射的原理

VSCode Debugger 里打的断点是怎么在网页里生效的？

1. VSCode 会记录你在哪个文件哪行打了个断点，然后会把这个信息发给 Chrome。
   这是一个形如 `/Users/guang/code/foo/src/App.vue` 的`绝对路径`
   但是问题来了，我们本地打的断点是一个绝对路径，也就是包含 `${workspaceFolder}` 的路径，而网页里根本没有这个路径，那怎么断住的呢？
   这是因为有的文件是关联了 sourcemap 的，它会把文件路径映射到源码路径
   如果映射到的源码路径直接就是本地的文件路径，那断点就生效了。
   ![alt text](image-2.png)
   vite 的项目，sourcemap 都是这种绝对路径，所以断点直接就生效了。
   但是 webpack 的项目，sourcemap 到的路径不是绝对路径，而是这种：
   ![alt text](image-3.png)
   本地打的断点都是绝对路径，而 sourcemap 到的路径不是绝对路径，根本打不上呀！
   `所以 VSCode Chrome Debugger 支持了 sourceMapPathOverrides 的配置，让打断点的文件路径和 sourcemap 之后的文件路径对上`
   如果映射之后路径开头多了几层目录，那就要配置下 webRoot。

## 用 VSCode 调试 React 项目

1. 需要 build 出带有 sourcemap 的 react 包
2. 下载之后 reset 到这个 commit：

   ```bash
   git reset --hard 80f3d88190c07c2da11b5cac58a44c3b90fbc296
   ```

3. 找到 rollup 的配置，添加一行 sourcemap: true，让 rollup 在构建时产生 sourcemap：
4. 找出没有生成 sourcemap 的那几个插件注释掉

## 文件只读，可能是本地不存在这个文件

把调试的文件 sourcemap 到的路径`映射`到本地的文件，这样调试的代码就不再只读了：

## json 与 jsonc

scope 指定为 json 和 jsonc，这是因为 json 文件对应两种语言：

jsonc 是 json with comments，带注释的 json，因为 json 语法是不支持注释的，而我们又想在 json 文件里加一些注释，所以平时都是用 jsonc 的类型。

## 调试 npm scripts

我们也可以用 npx 来跑，比如 npx xx，它的作用就是执行 node_modules/.bin 下的本地命令，如果没有的话会从 npm 下载然后执行。
npm scripts 本质上还是用 node 来跑这些 script 代码，所以调试他们和调试其他 node 代码没啥区别。
`命令行工具都是在 package.json 中声明一个 bin 字段，然后 install 之后就会放到 node_modules/.bin 下`
`可以 node node_modules/.bin/xx 来跑，可以 npx xx 来跑，最常用的还是 npm scripts，通过 npm run xx 来跑。`

1. 也就是可以这样跑：
   在 .vscode/launch.json 的调试文件里，选择 node 的 launch program：
   用 node 执行 node_modules/.bin 下的文件，传入参数即可：

```json
"configurations": [
  {
    "type": "node",
    "request": "launch",
    "name": "Launch Program",
    "skipFiles": ["<node_internals>/**"],
    "program": "${workspaceFolder}/node_modules/.bin/xx.js",
    "args": ["--arg1", "value1"]
  }
]
```

![alt text](image-4.png)

2. 其实还有更简单的方式，VSCode Debugger 对 npm scripts 调试的场景做了封装，可以直接选择 npm 类型的调试配置：

```json
"configurations": [
  {
    "type": "pwa-node",
    "request": "launch",
    "runtimeExecutable": "npm",
    "runtimeArgs": ["run-script", "foo"],
  }
]
```

## 命令行工具的两种调试方式（以 ESLint 源码调试为例）

这些命令行工具都提供了两种入口：**命令行和 api。**

1. 命令行的方式调试 ESLint 源码
   探究 fix 的原理
   `npx eslint ./index.js --fix ` 做了什么？

   ```json
   {
     "type": "node",
     "name": "eslint 调试",
     "program": "${workspaceFolder}/node_modules/.bin/eslint",
     "args": ["./index.js", "--fix"],
     "skipFiles": ["<node_internals>/**"],
     "console": "integratedTerminal",
     "cwd": "${workspaceFolder}",
     "request": "launch"
   }
   ```

   在 .bin 下找到 eslint 的文件打个断点

   eslint 的实现原理：
   ![alt text](image-5.png)
   lint 的实现是基于 AST，调用 rule 来做的检查。
   fix 的实现就是字符串的替换，多个 fix 有冲突的话会循环多次修复，默认修复 10 次还没修复完就终止。

2. api 的方式调试 ESLint 源码
   ESLint 会创建 ESLint 实例，然后调用 lintText 方法来对代码 lint。

   ESLint 源码的调试还是相对简单，因为没有经过编译，如果做了编译的话，那就需要 sourcemap 了

## 有时候我们需要修改 node_modules 下的一些代码，但是 node_modules 不会提交到 git 仓库，改动保存不下来，怎么办呢？

https://github.com/ds300/patch-package#readme
**这时候可以用 patch-package 这个工具。**
当我们需要对 node_modules 下的代码做改动的时候，可以通过 patch-package xxx 生成 patches 文件，它可以被提交到 git 仓库，然后再拉下来的代码就可以通过 patch-package 来应用改动。

1. 使用:

```bash
# 对node_modules/some-package/brokenFile.js做改动
vim node_modules/some-package/brokenFile.js

# 生成patch文件
npx patch-package some-package

# 提交patch文件
git add patches/some-package+3.14.15.patch
git commit -m "fix brokenFile.js in some-package"

# 把项目拉下来的时候，执行下 npx patch-package 就会应用这次改动
# 可以把它配到 postintsll 里，每次安装完依赖自动跑。
# 这样能保证每次拉取下来的代码都包含了对 node_modules 的改动。
```

2. 调试 patch-package 源码
   它默认就是有 sourcemap 的，只不过是 base64 的方式内联的(ts 配置 inlinsourcemap: true)
   ![alt text](image-6.png)

   探究它的实现原理要分为两各方面：
   **一个是 patches 文件怎么生成的，一个是 patches 文件怎么被应用的。**

   1. patches 文件怎么生成的(generate)
      看 patches 文件的内容就能看出来这是 git 的 diff：
      patch-package 就是依赖 git diff 实现的 patches 文件生成
      你可以先对 node_modules 下的某个包做下改动，然后执行 node ./dist/index xxx 来生成 patches 文件

      `好巧妙的方法！`：
      **在临时目录生成 package.json，下载依赖，生成一个 commit，然后把改动的代码复制过去，两者做 gif diff，就可以生成 patches 文件**

   2. patches 如何被应用的(apply)
      patch-package 自己实现了它的 parse，解析 patch 文件，拿到对什么文件的哪些行做什么修改的信息，之后根据不同做类型做不同的文件操作就可以了
      如果是 pnpm，那 patch-package 不支持，这时候用内置的 pnpm patch 命令就好了。
      `pnpm 内置 patch、patch-commit 命令，作用和这个 patch-package 包一样`

## 调试 Babel 源码

Babel 是一个 JS 的编译器，用于把高版本语法的代码转成低版本的，并且添加 polyfill。
它有很多插件，插件还进一步封装成了预设（preset），开箱即用。
此外，我们还可以写 Babel 插件来完成一些特定的代码转换。
`@babel/parser、@babel/traverse(transform)、@babel/generator `

```js
const parser = require('@babel/parser')
const traverse = require('@babel/traverse').default
const generate = require('@babel/generator').default

const source = `
    (async function() {
        console.log('hello guangguang');
    })();
`

const ast = parser.parse(source)

traverse(ast, {
  StringLiteral(path) {
    path.node.value = path.node.value.replace('guangguang', 'dongdong')
  }
})

const { code, map } = generate(ast, {
  sourceMaps: true
})

console.log(code)
console.log(JSON.stringify(map))
```

1. 怎么调试最初的源码呢
   sourcemap！
   但是你去 node_modules 下看下这些包，会发现它们已经有 sourcemap 了，而且也关联了：
   那为什么调试的时候调试的不是源码呢？
   这是因为 VSCode 的一个默认配置导致 sourcemap 不会生效。
   `resolveSourceMapLocations`!
   VSCode Node Debugger 默认不会查找 node_modules 下的 sourcemap。
2. 虽然调试的是源码的 ts 了，但是路径是 node_modules 包下的
   我们可以把 babel 项目下下来和测试项目放在一个 workspace 下，`然后去 node_modules 下手动替换下 sourcemap 的 sources 路径`，换成本地的路径，这样就可以调试 babel 源码了。
   然后在新的 workspace 创建个调试配置，这时目录改了，要指定下 cwd

## 实战案例：调试 Vite 源码

问题：vite 跑 dev server 的过程都执行了什么逻辑

1. debug npm script

```json
{
  "name": "Launch via NPM",
  "type": "node",
  "request": "launch",
  "runtimeExecutable": "npm",
  "runtimeArgs": ["run-script", "dev"],
  "console": "integratedTerminal",
  "skipFiles": ["<node_internals>/**"]
}
```

可以调用栈看到这部分代码是 node_modules/vite/dist/node 下的
去 node_modules 下看了下，并没有 sourcemap：
那去哪里找 sourcemap 呢？
这时就只能通过源码 build 了。
vite 是用 rollup 打包的，每个包下都有个 rollup.config.ts 文件，
搜一下 sourcemap，会找到一个 createNodeConfig 的函数，这里就是配置 node 部分的代码是否生成 sourcemap 的地方：
vite 编译会生成三部分代码，一部分是浏览器里的，也就是 client 目录下的，一部分是 node 跑的，是 esm 的模块，还有一部分是 node 跑的 cjs 的模块。

## 实战案例：调试 TypeScript 源码 (ts 的一些特性，可以通过源码找答案)

https://juejin.cn/book/7070324244772716556/section/7137086397147512840?utm_source=profile_book

```ts
type Test<T> = T extends number ? 1 : 2
type res = Test<any>

// 为什么res是1|2呢？
// 这就要从源码找答案了!
```

1. 下载源码
   lib 目录下有 `tsc.js` 和 `typescript.js`，这两个分别是 ts 的命令行和 api 的入口
   但是，这些是编译以后的 js 代码，源码在 src 下，是用 ts 写的。
   怎么把编译后的 js 代码和 ts 源码关联起来呢？ sourcemap！
2. 打断点
   ts 代码太多了，不知道哪些是解析类型的逻辑，在哪里打断点比较好。
   `这种情况还是用 api 的方式调试比较好。`

   ```js
   const ts = require('./built/local/typescript')

   const filename = './input.ts'
   const program = ts.createProgram([filename], {
     allowJs: false,
     strictNullChecks: true
   })

   const ast = program.getSourceFile(filename) // 想知道代码哪部分是什么 AST 可以通过 astexplorer.net 来查看
   const typeChecker = program.getTypeChecker()

   function visitNode(node) {
     if (node.kind === ts.SyntaxKind.TypeAliasDeclaration && node.name.escapedText === 'res') {
       const type = typeChecker.getTypeFromTypeNode(node.name)

       console.log(type)
     }

     node.forEachChild(child => visitNode(child))
   }

   visitNode(ast)
   ```

   ![alt text](image-7.png)
   通过运用 ts 的 api 去调试 ts 会比直接调试源码简单，如果条件类型左边是 any 时会作为联合类型处理

## 如何通过变量写出更灵活的调试配置？

vscode 的 debug 配置变量支持变量

- input 变量，可以让**用户输入或者选择**，通过 ${input:xxx} 语法
  例如："${input:port}" -> 输入一个端口号
- env 变量，可以读取**环境变量**值，通过 ${env:xxx} 语法
  例如："${env:HOME}" -> 取 HOME 环境变量
- config 变量，取 **vscode 的配置**，通过 ${config:xxx} 语法
  例如，可以通过 ${config:launch.nodePath} 来取 launch 配置里的 nodePath
- command 变量，可以**读取命令执行结果**，通过 ${command: xxx} 语法
  例如："${command:extension.pickNodeProcess}" -> 选择一个 node 进程
- 内置变量，可以取当前文件、目录等信息，通过 ${xxx} 语法
  例如：${file}、${workspaceFolder}、${relativeFile}、${fileBasename}、${fileDirname}、${fileExtname}、${cwd}

  灵活运用这些变量，可以让调试配置更灵活。

## 如何灵活的调试 Jest 测试用例

- 如果只是想跑某个测试文件的用例: jest 后面加上要跑的测试文件的就行
- 如果想跑某个测试用例: jest 后面加上 -t 参数，然后加上要跑的测试用例的名字。-t 是 --testNamePattern 的缩写，可以指定要跑的用例名的正则。

1. 跑 jest 怎么跑呢？
   其实我们跑 jest 最终执行的是 node_modules/jest/bin/jest.js 这个文件，所以调试的时候就直接用 node 跑这个文件，传入参数就行。
   还要指定日志输出位置为内置的终端，也就是 console 为 integratedTerminal。
   ![alt text](image-8.png)
   但你会发现它跑了多个 woker 进程，每个用例一个，这是 jest 优化性能的方式。
   但调试的时候可以不用这种优化，直接在主进程跑就行。
   可以加个 -i 的参数：
   -i 是 `--runInBand` 的缩写，这个参数的意思是不再用 worker 进程并行跑测试用例，而是在当前进程串行跑：
   但这样每调试一个用例都得改下配置也太麻烦了，能不能我打开哪个文件，就跑哪个文件的用例呢？
   可以的。
   VSCode 调试配置支持变量，比如 ${file} 就代表当前文件。
   这样就可以打开哪个调试那个了。
   那想指定具体的测试用例呢？
   **vscode 还支持输入类型的变量。**
2. jest-runner VSCode 插件 原理
   它会在每个测试用例旁边加一个运行和调试的按钮：
   点击不同位置的 debug，就是获取文件名和用例名传入调试配置
   相比第三方的 Jest 插件，自己写调试配置明显灵活、强大的多(例如正则)。

## 7 种断点(breakpoint)

这里只记录 3 种 chrome devtools 的 source 面板的

- dom 断点
  ![alt text](image-9.png)
  有三种类型，`子树修改的时候断住`、属性修改的时候断住、节点删除的时候断住。
  这时候你会发现代码在修改 DOM 的地方断住了，**这就是 React 源码里最终操作 DOM 的地方，看下调用栈就知道 setState 之后是如何更新 DOM 的了。**
- Event Listener 断点 (事件断点)
  之前我们想调试事件发生之后的处理逻辑，需要找到事件监听器，然后打个断点。
  打开 sources 面板，就可以找到事件断点，有各种类型的事件：
  ![alt text](image-11.png)
- url 请求断点
  当你想在某个请求发送的时候断住，但你不知道在哪里发的，这时候就可以用 url 请求断点
  ![alt text](image-12.png)
  不输入内容就是在任何请求处断住，你可以可以输入内容，那会在 url 包含该内容的请求处断住：

## 实战案例：调试 Ant Design 组件源码

1. 在组件里打个断点，代码会在这里断住
2. 在调用栈里找到`renderWithHooks`，这是 react 源码里调用函数组件的地方
3. 所有函数组件都是在这里被调用的，`而 antd 的组件也全部是函数组件，那么我们在这里加个断点，打名字为 Button 的函数组件被调用的时候断住不就行了？`
   右键选择添加条件断点：当组件名字包含 Button 的时候才断住。
   ![alt text](image-13.png)
4. InternalButton 在这里断住了。
   这个 InternalButton 就是 antd 里的 Button 组件。
5. step into 进入函数内部
   你会发现这确实是 Button 组件的源码，但却是被编译后的，比如 jsx 都被编译成了 React.createElement：
   这样是可以调试 Button 组件源码的，但是比较别扭。
   那能不能直接调试 Button 组件对应的 tsx 源码呢？
   可以的，这就要用到 sourcemap 了。
   ```bash
   git clone --depth=1 --single-branch git@github.com:ant-design/ant-design.git
   ```
   --single-branch 是下载单个分支， --depth=1 是下载单个 commit， 这样速度会快几十倍
6. 但你会发现 package.json 中有 build 命令，有 dist 命令，该执行哪个呢？
   这个就需要了解下 antd 的几种入口了。
   去 react 项目的 node_modules 下，找到 antd 的 package.json 看一下，你会发现它有三种入口：
   ![alt text](image-14.png)
   main 是 commonjs 的入口，也就是 require('antd') 的时候会走这个。
   module 是 esm 的入口，也就是 import xx from 'antd' 的时候会走这个。
   unpkg 是 UMD 的入口，也就是通过 script 标签引入的时候或者 commonjs 的方式等都可以用。
   分别对应了 lib、es、dist 的目录。
   `所以 antd 项目里的 dist 命令就是单独生成 UMD 代码的，而 build 命令是生成这三种代码。`
7. 那直接用 dist 入口的代码就能调试源码了么？
   把引入组件的地方换成 dist 目录下，也就是用 UMD 形式的入口。
   重新跑调试：
   你会发现代码确实比之前更像源码了。
8. 也就是没了 babel runtime 的代码，这明显是源码了。
   `但是依然还是 React.createElement，而不是 jsx，也没有 ts 的代码。`说明它还不是最初的源码。
   为什么会出现这种既是源码又不是源码的情况呢？
   **缺少 loader 的 sourcemap!**

   因为它的编译流程是这样的：
   ![alt text](image-15.png)
   Button.tsx --tsc -> Button.js --babel -> Button.js --webpack -> bundle.js
   tsc 和 babel 的编译都会生成 sourcemap，而 webpack 也会生成一个 sourcemap
   webpack 的 sourcemap 默认只会根据最后一个 loader 的 sourcemap 来生成。
   所以想映射回最初的 tsx 源码，只要关联了每一级 loader 的 sourcemap 就可以了(`也就是tsc和babel的sourcemap`)。而这个是可以配置的，就是 devtool。

9. 改编译配置，sourcemap 直接顶配
   antd 的编译工具链在 @ant-design/tools 这个包里，从 antd/node_modules/@antd-design/tools/lib/getWebpackConfig.js 就可以找到 webpack 的配置
10. 改 webpack 的 devtool：
    搜一下 ts-loader，你就会看到这段配置：
    ![alt text](image-16.png)
    确实就像我们分析的，tsx 会经过 ts-loader 和 babel-loader 的处理。
    搜一下 devtool，你会发现**它的配置是 source-map：**
    这就是 antd 虽然有 sourcemap，但是关联不到 tsx 源码的原因(没带`module`)。
    那我们给它改一下：
    把 devtool 改为 cheap-`module`-source-map，关联 loader 的 sourcemap
11. 改 babel：并且改一下 babel 配置，设置 sourceMap 为 true，让它生成 sourcemap
12. 改 ts：ts 也同样要生成 sourcemap，不过那个是在根目录的 tsconfig.json 里改
13. 重新 build，把它复制到 react 项目的 node_modules/antd/dist 下，覆盖之前的
14. 清一下 babel-loader 的缓存，删除整个 .cache 目录
15. 有的同学可能会担心 node_modules 下的改动保存不下来。
    这个也不是问题，可以执行下 npx patch-package antd，会生成这样一个 patch 文件:

## 实战案例：调试 ElementUI 组件源码

1. 我怎么知道在哪里打断点呢？
   我们可以知道的是，这个 button 会处理点击事件，但是却不知道事件处理函数的代码在什么地方。
   `这种情况可以加一个事件断点：`
   当你知道这个组件处理了什么事件，但却不知道事件处理函数在哪的时候就可以用事件断点。
2. methods、computed、props，这明显是源码里的了。但你再往上走两步，会发现又不是最初的源码：
   template 变成了 render 函数，而且还有其他组件的代码，`这明显是被编译打包之后的代码。`
   这是一个把所有组件代码编译后打包到一起的文件。
   `要用到 sourcemap 了`
3. 下载源码
   进入 element 目录，安装依赖，你会遇到一个前端经常头疼的问题，`node-sass 安装报错了`
4. 构建
   npx webpack --config build/webpack.common.js
   我们的目标是生成带有 source-map 的代码，所以要改下配置
   devtool 为 cheap-module-source-map。module 是把中间 loader 产生的 sourcemap 也给合并到最终的 sourcemap 里，这样才能直接映射到源码。
5. 覆盖 node_modules

## 如何用 Performance 工具分析并优化性能

Performance 工具最重要的是分析`主线程的 Event Loop`，分析每个 Task 的耗时、调用栈等信息
每个函数的耗时也都显示在左侧，右侧有源码地址，点击就可以跳到 Sources 对应的代码。
![alt text](image-17.png)
因为渲染和 JS 执行都在主线程，在一个 Event Loop 中，会相互阻塞，如果 JS 有长时间执行的 Task，就会阻塞渲染，导致页面卡顿。
**所以，性能分析主要的目的是找到 long task，之后消除它。**
在 Performance 中宽度代表时间，`超过 50ms 就被认为是 Long Task，会被标红`。因为如果 16.7 ms 渲染一帧，那 50ms 就跨了 3、4 帧了。

## 面试必问的异步顺序问题，用 Performance 轻松理清

## Web Vitals 如何计算？如何衡量网页性能？

**3 个作为核心指标，也就是 LCP、FID、CLS，分别用来衡量加载性能、交互性能、视觉稳定性。**

LCP 是最大内容渲染时间，代表页面已经完成了主要内容的渲染，这个指标可以用来衡量加载到渲染的性能。（FMP 是有意义的渲染，那个比较难定义）

FID 是首次输入延迟，衡量页面内容首次渲染（FCP）之后，到可交互（TTI）的这段时间内，用户点击按钮或者输入内容到页面响应的时间。是从用户交互角度衡量页面性能的指标。

CLS 是累积布局偏移，布局稳定性，是能反应用户体验的一个指标。

## console.log 会导致内存泄漏？用 Performance 和 Memory 分析下（上）

1. 勾选 Memory，然后开始录制，点击 3 次按钮，再执行一次 GC：
2. performance 工具排查
3. 通过 performance.memory.totalJSHeapSize 是可以拿到堆内存大小
   **不打开 devtools 的时候，console.log 不会内存泄漏**

## Chrome DevTools 小功能集锦

1. **重新发送一次 XHR 请求**
   选中一个请求，右键，Replay XHR
2. **请求定位到源码**
   你想知道某个请求是在哪里发的，可以打开 Network 面板，在每个网络请求的 initiator 部分可以看到发请求代码的调用栈，点击可以快速定位到对应代码。
3. **元素定位到创建的源码**
   可以通过 Elements 面板选中某个元素，点击 Stack Trace，就会展示出元素创建流程的调用栈
   这个功能是实验性的，需要手动开启下：在 settings 的 expriments 功能里，勾选 “Capture node creation stacks”。
4. group by folder
5. Network 自定义展示列
   - Network 是可以修改展示的列的，比如我勾选 `Cookie 和 Set-Cookie：`
     Set-Cookie 的意思就是这个请求的响应会设置几个 cookie。
     Cookies 的意思就是说这个请求会携带几个 cookie：
   - 除此以外，你还可以自定义展示的响应 header：我比较常用的是 `Cache Control`：
   - 右边有个 waterfall，默认是展示请求的时间，但我觉得这个没啥用，我更喜欢看请求响应的耗时，所以我会把它换成` total duration：`，这样 waterfall 展示的就是耗时了。
6. 代码手动关联 sourcemap
   sources 面板可以右键点击 add soruce map，输入 sourcemap 的 url 就可以关联上 sourcemap，`当调试线上代码的时候可以用这种方式关联源码。`
7. **filter**
   输入一个 - 就会提示所有的过滤器，但是这个减号之后要去掉，它是非的意思
   ![alt text](image-18.png)
8. developer resources
   sourcemap 文件为啥在 Network 里看不到呢
   其实这个被 Network 过滤掉了，想看到这些文件的请求在另一个地方：
   ![alt text](image-19.png)
9. **remove event listeners**
   element 面板选中元素可以看到这个元素和它的父元素的所有事件监听器：
   可以手动 remove。
   **比如你想看下拉菜单的样式，但是鼠标一移开就消失了**
   **这时候你可以删掉这个按钮的 mouseleave 事件的监听器：**
   ![alt text](image-20.png)

## 通过调试技术，理清 b 站视频播放速度很快的原因

range 是一个 http 请求头，用来指定请求资源的范围，比如 range: bytes=0-1000 就是请求资源的前 1000 个字节。无法其返回 206 partial content 状态码。如果请求的范围不合法，就会返回 416 Range Not Satisfiable 状态码。

知道了为什么 b 站视频播放的那么快了：
`视频分成多少段是提前就确定的，大概 600kB 一个片段，m4s 分段存储视频，通过 range 请求动态下载某个视频片段。然后通过 SourceBuffer 来动态播放这个片段。`

## 用 Debugger 你能读懂各种源码

`用 Performance 工具看到代码执行的全貌，然后用 Debugger 来深入每一条代码执行路径的细节。`

## Chrome DevTools 的原理

`Chrome DevTools 被设计成了和 Chrome 分离的架构，两者之间通过 WebSocket 通信，设计了专门的通信协议： Chrome DevTools Protocol。`
这样只要实现对接 CDP 协议的 ws 服务端(backed)，就可以用 Chrome DevTools 来调试，**所以 Chrome DevTools 可以用来调试浏览器的网页、Node.js，Electron、安卓手机的网页 等**

## 线上报错如何直接调试本地源码？

`通过 sourcemap，我们可以调试线上报错的时候直接对应到本地源码来断点调试。`
VSCode Debugger 异常断点断住的时候右键 chrome devtools 的 sources 面板，手动 add source map。不过这种方式是一次性的。
`解决一次性的方式`：在 Network 面板中，右键报错文件的请求，选择 `Override content`，在 Sources 面板中的 Overrides 就可以修改文本内容，刷新重新请求就会走该文件。

## 远程调试移动端网页的原理(实现 backed)

我们可以通过 chrome、safari 调试移动端的网页，原理就是开启调试模式之后，**可以通过 CDP server 和 client 进行通信，从而实现调试。**
可以通过 USB 调试，是因为 adb 做了端口转发，也可以通过 wifi 调试，这种就需要自己实现一个 ws 服务做中转了。

## 手动对接 CDP 来理解 Chrome DevTools 实现原理(实现 fronted)

https://juejin.cn/book/7070324244772716556/section/7137086339941400612?scrollMenuIndex=0

**cdp 协议文档：https://chromedevtools.github.io/devtools-protocol/**

CDP 协议可以在 Protocol Monitor 里看到，分成了多个 Domain，每个 Domain 下有很多 Method 和 Event。
![alt text](image-21.png)

- Method 就是 frontend 向 backend 请求数据，backend 给它返回相应的数据
- Event 就是 backend 推送给 frontend 的一些数据。

Chrome DevTools 里展示的所有内容都是从 backend 那里拿到的，他只是一个展示和交互的 UI 而已。
这个 UI 是可以换的，比如我们可以用 VSCode Debugger 对接 CDP 调试协议来调试网页。
有很多用 Chrome DevTools 调试别的目标的工具，而前端领域的跨端引擎、小程序引擎也是通过这种方式实现的。
跨端引擎要自己实现 CDP 协议的对接，而小程序引擎简单一些，本来就有了 CDP backend，对接到 frontend 即可。

可以自己实现 CDP backend，当然也可以实现 frontend，但**自己对接 WebSocket 和 CDP 协议还是挺麻烦的，可以直接用 chrome-remote-interface 这个包。**

## 基于 CDP 实现 Puppeteer

**puppeteer 是基于 Chrome DevTools Protocol 实现的，会以调试模式跑一个 chromium 的实例，然后通过 WebSocket 连接上它，之后通过 CDP 协议来远程控制**。
**我们写的脚本最终都会转成 CDP 协议来发送给 Chrome 浏览器，这就是它的实现原理。**

## 把 Puppeteer 融入调试流程来做自动化

https://juejin.cn/book/7070324244772716556/section/7182106489870843941?scrollMenuIndex=0

一个不错的 idea

这俩其实完全可以结合在一起用，因为他们都是基于 CDP，会启动一个调试模式的浏览器。
**等 puppeteer 跑起 chrome 之后，vscode debugger 再 attach 上它来调试**

融合在一起之后，你可以写 puppeteer 脚本来自动化一些流程，比如自动登录、自动填写表单等，这个过程还可以断点调试，断点释放之后再执行后续自动化脚本
`把 puppeteer 融入调试流程，调试体验爽翻了！`

---

更新：

- Chrome DevTools 是一个前后端项目，后端 ws，接口 cdp，前端普通的 html。

- 跨端引擎的 renderView 不同，dom 结构相同。
  看到自己模拟 DOM 信息这部分，不知道你是否会想到跨端引擎呢。
  跨端引擎就是通过前端的技术来描述界面（比如也是通过 DOM），实际上用安卓和 IOS 的原生组件来做渲染。
  它的调试工具也是需要显示 DOM 树的信息的，但是因为并不是网页，所以不能直接用 Chrome DevTools。
  那如何用 Chrome DevTools 来调试跨端引擎呢？
  看完上面两个案例，相信你就会有答案了。只要对接了 CDP，自己实现一个 backend，把 DOM 树的信息，通过 CDP 的格式传给 frontend 就可以了。
  自定义的调试工具基本都是前端部分集成下 Chrome DevTools frontend，后端部分实现下对接 CDP 的 ws 服务来实现的。
  这节我们探究了下 Chrome DevTools 的实现原理，它只是一个对接了 CDP 的 UI，完全可以用来调试别的 target，只要实现对应的 CDP backend 即可。

- CDP 协议可以在 Protocol Monitor 里看到，分成了多个 Domain，每个 Domain 下有很多 Method 和 Event。
  有很多用 Chrome DevTools 调试别的目标的工具，而前端领域的跨端引擎、小程序引擎也是通过这种方式实现的。
  跨端引擎要自己实现 CDP 协议的对接，而小程序引擎简单一些，本来就有了 CDP backend，对接到 frontend 即可
  如果某一天，你也要做一个网页分析工具，是不是也可以通过 CDP 的方式来获取一些网页运行数据做分析呢？
  所有 Chrome DevTools 的数据，你通过 CDP 都是能拿到的，能做的事情有很多。
