1. 函数中的布尔参数
   函数中的布尔参数可能是浪费大量时间的原因，如果使用不当，可能导致代码可读性低。它们有时被认为是一种反模式，因为它们增加了认知负荷，降低了共享代码的可维护性。幸运的是，使用纯 JavaScript 选项对象很容易发现和修复它们。

```JS
// Real quick: Is this valid or invalid?
input.setInvalid(false);


缓解这个问题
// Ok, so reload but not immediately
results.reload({ immediate: false });
// Create a new user without administrator privileges
const user = new User({ isAdministrator: false });
```

2. forEach 里不要用异步任务
   使用 Array.prototype.forEach ()的异步回调将导致其余的代码执行，并且不会等待异步操作。
   唯一可行的解决方案是使用 Promise.all ()
3. 关于合并数组
   扩展操作符比 concat 性能更好
   然而，Array.prototype.concat ()能够比 spread 操作符更好地处理非数组值

```JS
const a = [1, 2, 3];
const b = true;
const c = 'hi';

const spreadAb = [...a, ...b]; // Error: b is not iterable
const spreadAc = [...a, ...c]; // [1, 2, 3, 'h', 'i'], wrong result
// You should use [...a, b] and [...a, c] instead

const concatAb = [].concat(a, b); // [1, 2, 3, true]
const concatAb = [].concat(a, c); // [1, 2, 3, 'hi']
```

只要知道输入是数组，就使用扩展运算符(...) ，因为它性能更好，易于阅读和理解。当您不能确定一个或多个输入并且不希望添加额外的检查时，使用 ary.prototype.concat () ，因为它能更优雅地处理这些情况。

4. deep freeze
   要使对象不可变，我们可以使用 Object.freeze () ，它将防止添加新属性，并在一定程度上防止删除和更改现有属性
   它执行的是浅冻结
5. encodeURI ()和 encodeURIComponent () 有什么区别
   encodeURI()不会对本身属于 URI 的特殊字符进行编码，例如**冒号、正斜杠、问号和井字号**；
   而 encodeURIComponent()则会对它发现的任何非标准字符进行编码

```JS
const partOfURL = 'my-page#with,speci@l&/"characters"?';
const fullURL = 'https://my-website.com/my-page?query="a%b"&user=1';

encodeURIComponent(partOfURL); // Good, escapes special characters
// 'my-page%23with%2Cspeci%40l%26%2F%22characters%22%3F'

encodeURIComponent(fullURL);  // Bad, encoded URL is not valid :https://的`://` 全部转掉了
// 'https%3A%2F%2Fmy-website.com%2Fmy-page%3Fquery%3D%22a%25b%22%26user%3D1'
```

```JS
const partOfURL = 'my-page#with,speci@l&/"characters"?';
const fullURL = 'https://my-website.com/my-page?query="a%b"&user=1';

encodeURI(partOfURL); // Bad, does not escape all special characters  有些字符没转
// 'my-page#with,speci@l&/%22characters%22?'

encodeURI(fullURL);  // Good, encoded URL is valid
// 'https://my-website.com/my-page?query=%22this%25thing%22&user=1'
```

encodeURI:允许更多字符
encodeURIComponent:更加严格

5. HTMLBodyElement
   ![继承关系](https://yari-demos.prod.mdn.mozit.cloud/zh-CN/docs/Web/API/HTMLBodyElement/_sample_.inheritance_diagram.html)

6. Hasinstance 允许我们定制 instanceof 运算符的行为

```JS
class PrimitiveNumber {
  static [Symbol.hasInstance] = x  => typeof x === 'number';
}
123 instanceof PrimitiveNumber; // true
```

7. 为了防止在新选项卡中打开的链接引起任何麻烦，我们应该总是将 rel = “ noopener noreferrer”属性添加到所有目标 = “ \_ blank”链接中。

```HTML
<!-- Bad: susceptible to tabnabbing -->
<a href="https://externalresource.com/some-page" target="_blank">
  External resource
</a>

<!-- Good: new tab cannot cause problems  -->
<a
  href="https://externalresource.com/some-page"
  target="_blank"
  rel="noopener noreferrer"
>
  External resource
</a>
```

**標籤釣魚 Tabnabbing**
当链接到我们网站的外部资源时，我们使用 target = “ \_ blank”在新的标签页或窗口中打开链接页面
但是有一个安全风险我们应该意识到。新标签通过 **window.opener** 获得了对链接页面(即我们的网站)的有限访问权限，然后可以通过 window.opener.location 更改链接页面的 URL(这被称为 **tabnabbing**)
window.opener:返回打开当前窗口的那个窗口的引用，例如：在 window A 中打开了 window B，B.opener 返回 A.

如果外部资源是不值得信任的，可能已经被黑客攻击，域名已经更换所有者多年等
为了防止在新选项卡中打开的链接引起任何麻烦，我们应该总是将 rel = “ noopener noreferrer”属性添加到所有目标 = “ \_ blank”链接中。
(即，不能获得我们网站 A 窗口的引用)

8. 如何在 JavaScript 中清空数组？
   将其长度设置为 0
   这个方法也非常快，并且具有处理 const 变量的额外好处。

```JS
let a = [1, 2, 3, 4];
a.length = 0;
```

9. 如何让搜索引擎抓取 AJAX 内容
   搜索引擎只抓取 example.com，不会理会井号，因此也就无法索引内容
   解决：Discourse 是一个论坛程序，严重依赖 Ajax，但是又必须让 Google 收录内容。它的解决方法就是放弃井号结构，**采用 History API**
