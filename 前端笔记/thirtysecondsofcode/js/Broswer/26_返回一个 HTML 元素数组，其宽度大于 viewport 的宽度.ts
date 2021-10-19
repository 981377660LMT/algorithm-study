getElementsBiggerThanViewport() // <div id="ultra-wide-item" />

function getElementsBiggerThanViewport() {
  const docWidth = document.documentElement.offsetWidth
  return Array.from<HTMLElement>(document.querySelectorAll('*')).filter(
    el => el.offsetWidth > docWidth
  )
}
