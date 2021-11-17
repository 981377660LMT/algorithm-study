1. HTML 事件处理程序：用户可能在元素刚出现就触发了事件，但此时 JS 代码可能还未加载完毕。其次，HTML 代码和 JavaScript 代码紧密耦合，不利于开发和维护，所以不推荐使用这种方法。
   DOM0 级事件处理程序：简单且兼容性好，但是在需要对一个元素设置多个事件处理程序时便显得孱弱。
   DOM2 级事件处理程序：可以轻易的设置多个事件处理程序，但是在删除事件处理程序时，传给 removeEventListener() 的参数必须与之前一致，且 IE9 以下不支持。

2. js 的各种位置，比如 clientHeight,scrollHeight,offsetHeight ,以及 scrollTop, offsetTop,clientTop 的区别？

   1. clientHeight：表示的是可视区域的高度，不包含 border 和滚动条

   2. offsetHeight：表示可视区域的高度，包含了 border 和滚动条

   3. scrollHeight：表示了所有区域的高度，包含了因为滚动被隐藏的部分。

   4. clientTop：表示边框 border 的厚度，在未指定的情况下一般为 0

   5. scrollTop：滚动后被隐藏的高度，获取对象相对于由 offsetParent 属性指定的父坐标(css 定位的元素或 body 元素)距离顶端的高度。

3. 图片懒加载三种方式
   1. offsetTop < clientHeight + scrollTop
   2. element.getBoundingClientRect().top < clientHeight
   3. IntersectionObserver 方式;intersectionRatio：目标元素的可见比例，即 intersectionRect 占 boundingClientRect 的比例，完全可见时为 1 ，完全不可见时小于等于 0
