1. 某主机的 IP 地址为 212.212.77.55，子网掩码为 255.255.252.0。若该主机向其所在子网发送`广播`分组，则目的地址可以是？
   由子网掩码可知前 22 位为`子网号`、后 10 位为`主机号`
   `主机号全部置为 1 即广播地址`
   212.212.（100011-01）.55 => 212.212. 100011-11.255
   在使用 TCP/IP 协议的网络中，主机标识段 host ID 为全 1 的 IP 地址为广播地址
2. 标签模板其实不是模板，而是函数调用的一种特殊形式。“标签”指的就是函数，紧跟在后面的模板字符串就是它的参数。

```JS
let a = 5; let b = 10; tag`Hello ${ a + b } world ${ a * b }`; // 等同于 tag(['Hello ', ' world ', ''], 15, 50);
```

3. 关于将 Promise.all 和 Promise.race 传入空数组的两段代码的输出结果说法正确的是：

```JS
Promise.all([]).then((res) => {
    console.log('all');
});
Promise.race([]).then((res) => {
    console.log('race');
});
```

all 会被输出，而 race 不会被输出
Promise.all([ ])中，数组为空数组，则立即决议为成功执行 resolve( )；
Promise.race([ ])中数组为空数组，就不会执行，永远挂起

4. 内存泄漏是 javascript 代码中必须尽量避免的，以下几段代码可能会引起内存泄漏的有（）

意外的全局变量

```JS
function getName() {
    name = 'javascript'
}
getName()
```

脱离 DOM 的引用
// 此时，仍旧存在一个全局的 #button 的引用
// elements 字典。button 元素仍旧在内存中，不能被 GC 回收。
`应该直接 elements.button=null`

```JS
const elements = {
    button: document.getElementById('button')
};
function removeButton() {
    document.body.removeChild(elements.button);
}
removeButton()
```

5. 身在乙方的小明去甲方对一网站做渗透测试，甲方客户告诉小明，目标站点由 wordpress 搭建，请问下面操作最为合适的是
   使用 wpscan 对网站进行扫描
6. 发生冲突的情况下，max-height 会覆盖掉 height, min-height 又会覆盖掉 max-height
7. 关于网络请求延迟增大的问题，以下哪些描述是正确的()
   使用 ping 来测试 TCP 端口是不是可以正常连接 (x) `ping 是网络层的，tcp 是传输层的`
   使用 tcpdump 抓包分析网络请求 (√)
   使用 strace 观察进程的网络系统调用 (√)
   使用 Wireshark 分析网络报文的收发情况 (√)
8. JSON.stringify 方法，返回值 res 分别是

   ```JS
   const fn = function(){}
   const res = JSON.stringify(fn)
   const num = 123
   const res = JSON.stringify(num)
   const res = JSON.stringify(NaN)
   const b = true
   const res = JSON.stringify(b)

   undefined、'123'、'null'、'true'
   ```

   1. 转换值如果有 toJSON() 方法，该方法定义什么值将被序列化。
   2. 函数、undefined 被单独转换时，会返回 undefined，如 JSON.stringify(function(){}) or JSON.stringify(undefined).
   3. 对包含循环引用的对象（对象之间相互引用，形成无限循环）执行此方\*\*\*抛出错误。
   4. NaN 和 Infinity 格式的数值及 null 都会被当做 null。
   5. undefined、任意的函数以及 symbol 值，在序列化过程中会被忽略（出现在非数组对象的属性值中时）或者被转换成 null（出现在数组中时）。

9. SASS/SCSS 依赖 Ruby 或 Node.js 编译环境，因此需要编译后才能在浏览器使用；而 less 本身由 JavaScript 实现，可以在浏览器中完成编译并可直接使用
10. 若一个哈夫曼树有 N 个叶子节点，则其节点总数为 2N-1

11. 分页存储管理将进程的逻辑地址空间分成若干个页，并为各页加以编号，从 0 开始，若某一计算机主存按字节编址，逻辑地址和物理地址都是 32 位，`页表项`大小为 4 字节，若使用一级页表的分页存储管理方式，逻辑地址结构为页号（20 位），`页内偏移量`（12 位），则页的大小是（4KB ）？页表最大占用（4MB ）？

地址长度为 32 位，其中 0~11 位为页内地址（即页内偏移量），2^12 即`每页大小为 4KB；`
同样地，12~31 位为页号，地址空间最多允许有 2^20 = 1M 页，又页表项 4 字节， 所以页表最大占用` 1M * 4 = 4MB`
