'use client'

import './style.scss'
import { RandContext } from '@/app/context'
import { useContext } from 'react'

const decorationImage = [
  "/dalle01.png", "/dalle02.png", "/dalle03.png",
  "/dalle04.png", "/dalle05.png", "/dalle06.png",
  "/dalle07.png"
]

const decorationText = [
  "Impressionist painting of Century Tower (AI Generated)",
  "Impressionist painting of UF building (AI Generated)",
  "Impressionist painting of UF building (AI Generated)",
  "Impressionist painting of UF stadium (AI Generated)",
  "Impressionist painting of Century Tower (AI Generated)",
  "Picture of Gator enjoying the view (AI Generated)",
  "Photo of Miniature Gator Skateboarder (AI Generated)",
]

const ImageBanner = () => {
  const random = useContext(RandContext)
  const src = decorationImage[random]
  const txt = decorationText[random]
  return (
    <div className="decoration">
      <div style={{overflow: "hidden"}}>
	  <img src={src} id="decoration--image"/>
      </div>
      <div id="decoration--meta">{txt}</div>
    </div>
  )
}

export default ImageBanner
