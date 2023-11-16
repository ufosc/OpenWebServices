import './imagebanner.scss'

const decorationImage = [
  "/dalle01.png", "/dalle02.png", "/dalle03.png",
  "/dalle04.png", "/dalle05.png",
]

const decorationText = [
  "Impressionist painting of Century Tower (AI Generated)",
  "Impressionist painting of UF building (AI Generated)",
  "Impressionist painting of UF building (AI Generated)",
  "Impressionist painting of UF stadium (AI Generated)",
  "Impressionist painting of Century Tower (AI Generated)"
]

const ImageBanner = () => {
  const random = Math.floor(Math.random() * decorationImage.length)
  const text = decorationText[random]
  const src = decorationImage[random]
  return (
    <div className="decoration">
      <div style={{overflow: "hidden"}}>
	  <img src={src} id="decoration--image"/>
      </div>
      <div id="decoration--meta">{text}</div>
    </div>
  )
}

export default ImageBanner
