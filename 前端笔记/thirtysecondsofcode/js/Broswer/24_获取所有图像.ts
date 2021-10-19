getImages(document, true) // ['image1.jpg', 'image2.png', 'image1.png', '...']
getImages(document, false) // ['image1.jpg', 'image2.png', '...']

function getImages(el: Document, includeDuplicates = false) {
  const images = Array.from(el.getElementsByTagName('img')).map(img => img.getAttribute('src'))
  return includeDuplicates ? images : [...new Set(images)]
}
