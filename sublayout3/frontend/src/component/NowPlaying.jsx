import React from 'react';
import Block from "./Block.jsx";
import Img1 from "../img/Screenshot 2024-05-12 at 3.35.59 PM.png"
import Img2 from "../img/Screenshot 2024-05-12 at 3.35.59 PM copy.png"
import Img3 from "../img/Screenshot 2024-05-12 at 3.35.59 PM copy 2.png"


const NowPlaying = () => {
	return (
		<>
			<Block style={{width: "100%", height: "96px", background: "#000", flexDirection: "row", justifyContent: "space-between"}}>
				<img src={Img3} alt="Im1"/>
				<img src={Img2} alt="Im1"/>
				<img src={Img1} alt="Im1"/>
			</Block>
			<Block
				style={{
					width: "calc(100% - 1rem)",
					margin: "0.5rem",
					height: "36px",
					background: "#1ed760",
					padding: "8px 24px",
					gap: "4px",
					flexDirection: "row",
					justifyContent: "end",
					color: "#000"
				}}>
				<svg data-encore-id="icon" role="presentation" aria-hidden="true"
				     style={{
					     width: 16,
				     }}>
					<path
						d="M14.5 8a6.468 6.468 0 0 0-1.3-3.9l1.2-.9C15.405 4.537 16 6.2 16 8c0 1.8-.595 3.463-1.6 4.8l-1.2-.9A6.468 6.468 0 0 0 14.5 8zM8 1.5a6.5 6.5 0 1 0 3.25 12.13.75.75 0 1 1 .75 1.3 8 8 0 1 1 0-13.86.75.75 0 1 1-.75 1.298A6.467 6.467 0 0 0 8 1.5z"></path>
					<path
						d="M11.259 8c0-.676-.228-1.296-.611-1.791l1.187-.918c.579.749.924 1.69.924 2.709a4.41 4.41 0 0 1-.925 2.709l-1.186-.918c.383-.495.61-1.115.61-1.791zM8.75 4.115l-4.139 2.39a1.727 1.727 0 0 0 0 2.99l4.139 2.39v-7.77z"></path>
				</svg>
				<h5 style={{marginTop: 2}}>Playing on CSC105 Pre-hackathon</h5>
			</Block>
		</>
	);
};

export default NowPlaying;
