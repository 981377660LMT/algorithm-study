1. 介绍下重绘和回流（Repaint & Reflow），以及如何进行优化
   1. 浏览器渲染机制
      浏览器采用流式布局模型（Flow Based Layout）
      浏览器会把 HTML 解析成 DOM，把 CSS 解析成 CSSOM，DOM 和 CSSOM 合并就产生了渲染树（Render Tree）。
      有了 RenderTree，我们就知道了所有节点的样式，然后计算他们在页面上的大小和位置，最后把节点绘制到页面上。
      由于浏览器使用流式布局，对 Render Tree 的计算通常只需要遍历一次就可以完成，但 table 及其内部元素除外，他们可能需要多次计算，通常要花 3 倍于同等元素的时间，这也是为什么要避免使用 table 布局的原因之一。
   2. 重绘：不会影响布局的变化例如 outline, visibility, color、background-color
   3. 回流：布局或者几何属性需要改变就称为回流
   4. 浏览器优化：现代浏览器大多都是通过队列机制来批量更新布局，浏览器会把修改操作放在队列中，至少一个浏览器刷新（即 16.6ms）才会清空队列，但当你获取布局信息的时候（offsetTop,width 等），队列中可能有会影响这些属性或方法返回值的操作，即使没有，浏览器也会强制清空队列，触发回流与重绘来确保返回正确的值。
   5. 减少回流与重绘
      使用 transform 替代 top
      使用 visibility 替换 display: none ，因为前者只会引起重绘，后者会引发回流（改变了布局
      CSS3 硬件加速（GPU 加速），使用 css3 硬件加速，可以让 transform、opacity、filters 这些动画不会引起回流重绘 。
      JS 批量修改 DOM
