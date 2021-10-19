/**
 *
 * @param str 将字符串复制到剪贴板。只有在用户操作(即在单击事件监听器中)的情况下才能工作。
 */
const copyToClipboard = (str: string) => {
  const el = document.createElement('textarea')
  el.value = str
  el.setAttribute('readonly', '')
  el.style.position = 'absolute'
  el.style.left = '-9999px'
  document.body.appendChild(el)

  // 使用 Selection.getRangeAt ()存储选定范围(如果有的话)。
  // const selected =
  //   document.getSelection()!.rangeCount > 0 ? document.getSelection()!.getRangeAt(0) : false
  el.select()
  document.execCommand('copy') // 将 < textarea > 的内容复制到剪贴板
  document.body.removeChild(el)
  // if (selected) {
  // 保存并恢复原始文档的选定内容
  //   document.getSelection()!.removeAllRanges()
  //   document.getSelection()!.addRange(selected) // Selection () . addRange ()恢复原来选择的范围(如果有的话)。
  // }
}

copyToClipboard('Lorem ipsum') // 'Lorem ipsum' copied to clipboard.
