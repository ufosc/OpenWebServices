'use client'

import './style.scss'
import { RandContext } from '@/app/context'
import { useContext, useState } from 'react'
import Image from 'next/image'
import { Loading } from '@carbon/react'

import dalle01 from "@/public/dalle01.png"
import dalle02 from "@/public/dalle02.png"
import dalle03 from "@/public/dalle03.png"
import dalle04 from "@/public/dalle04.png"
import dalle05 from "@/public/dalle05.png"
import dalle06 from "@/public/dalle06.png"
import dalle07 from "@/public/dalle07.png"

const decorationImage = [
  dalle01, dalle02, dalle03, dalle04,
  dalle05, dalle06, dalle07,
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
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const random = useContext(RandContext)
  const src = decorationImage[random]
  const txt = decorationText[random]
  return (
    <div className="decoration">
      <div style={{overflow: "hidden"}}>
	{
	  (isLoading) ?
	    (<Loading id="decoration--loading" withOverlay={false} />)
	    : null
	}
	<Image src={src} id="decoration--image" alt={txt}
	  onLoad={() => { setIsLoading(false) }}/>
      </div>
      <div id="decoration--meta">{txt}</div>
    </div>
  )
}

export default ImageBanner
