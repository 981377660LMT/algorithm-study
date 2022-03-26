// 图片懒加载

function lazyLoading(node: HTMLImageElement, src: string) {
  const intersectionObserver = new IntersectionObserver(entries => {
    if (entries[0].intersectionRatio > 0) return (node.src = src)
  })
  intersectionObserver.observe(node)
}

lazyLoading(document.querySelector('#img1')!, 'https://ok.com/bar.png')
