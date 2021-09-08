// 获取当前页面中所有 HTML tag 的 名字，以数组形式输出, 重复的标签不重复输出。（不考虑 iframe 和 shadowDOM）

function getAllTags() {
  const tags = Array.from(window.document.querySelectorAll('*')).map(tag => tag.tagName)
  return [...new Set(tags)]
}
