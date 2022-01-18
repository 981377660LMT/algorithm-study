import axios from 'axios'

const canvas = document.createElement('canvas')
const dataURL = canvas.toDataURL()
const form = new FormData()
const myimg = new Blob([dataURL])
form.append('my-img', myimg)

axios.post('url', form)
