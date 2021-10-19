// https://www.ruanyifeng.com/blog/2009/09/find_element_s_position_using_javascript.html

// 1.获取网页的大小
// 不包含滚动内容：clientHeight和clientWidth
// documentElement的clientHeight和clientWidth属性，就代表了网页的大小。
// 包含滚动内容：网页上的`每个元素`还有scrollHeight和scrollWidth属性，指`包含滚动条在内`的该元素的视觉面积。

// 2.获取网页元素的绝对位置:offsetTop迭代
// 每个元素都有offsetTop和offsetLeft属性，表示该元素的左上角与父容器（offsetParent对象）左上角的距离
// 所以，只需要将这两个值进行累加，就可以得到该元素的绝对坐标。
getVerticalOffset(document.querySelector('.my-element')!) // 120

function getVerticalOffset(el: HTMLElement) {
  let res = el.offsetTop
  while (el.offsetParent) {
    el = el.offsetParent as HTMLElement
    res += el.offsetTop
  }
  return res
}

// 获取网页元素的相对位置:offsetTop迭代+最后减去滚动距离

////////////////////////////////////////////////////////////////////
// 获取元素位置的快速方法 Element.getBoundingClientRect()
// 1.相对位置:相对于浏览器窗口（viewport）左上角的距离
document.querySelector('#app')!.getBoundingClientRect().top

// 2.绝对位置，包含滚动距离
document.querySelector('#app')!.getBoundingClientRect().top +
  (document.documentElement.scrollTop || document.body.scrollTop)
////////////////////////////////////////////////////////////////////////
// 检查页面底部是否可见
const bottomVisible = () =>
  document.documentElement.clientHeight +
    (document.documentElement.scrollTop || document.body.scrollTop) >=
  document.documentElement.scrollHeight
bottomVisible() // true

// 例子：
// 一开始进入页面：
// document.documentElement.clientHeight:702
// window.scrollY:0
// document.documentElement.scrollHeight:1856
