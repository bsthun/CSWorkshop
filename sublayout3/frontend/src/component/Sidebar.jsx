import React from 'react';
import Block from "./Block.jsx";
import style from "./Sidebar.module.scss";
import usePreview from "../hook/usePreview.js";

const Sidebar = () => {
	const sidebarUrl = usePreview("sidebar")
	return (
		<div style={{
			flex: 2,
			display: "flex",
			flexDirection: "column",
			alignItems: "stretch",
			justifyContent: "start",
			gap: "0.5rem",
			height: "100%",
		}}>
			<Block>
			<div className={style.navContainer}>
				<svg data-encore-id="icon" role="img" aria-hidden="true"
				     style={{
					     width: 24,
					     fill: "currentcolor"
				     }}
				     viewBox="0 0 24 24">
					<path
						d="M13.5 1.515a3 3 0 0 0-3 0L3 5.845a2 2 0 0 0-1 1.732V21a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1v-6h4v6a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V7.577a2 2 0 0 0-1-1.732l-7.5-4.33z"></path>
				</svg>
				<h4>Home</h4>
			</div>

			<div className={style.navContainer} style={{opacity: 0.8}}>
				<svg data-encore-id="icon" role="img" aria-hidden="true"
				     style={{
					     width: 24,
					     fill: "currentcolor"
				     }}
				     viewBox="0 0 24 24">
					<path
						d="M10.533 1.27893C5.35215 1.27893 1.12598 5.41887 1.12598 10.5579C1.12598 15.697 5.35215 19.8369 10.533 19.8369C12.767 19.8369 14.8235 19.0671 16.4402 17.7794L20.7929 22.132C21.1834 22.5226 21.8166 22.5226 22.2071 22.132C22.5976 21.7415 22.5976 21.1083 22.2071 20.7178L17.8634 16.3741C19.1616 14.7849 19.94 12.7634 19.94 10.5579C19.94 5.41887 15.7138 1.27893 10.533 1.27893ZM3.12598 10.5579C3.12598 6.55226 6.42768 3.27893 10.533 3.27893C14.6383 3.27893 17.94 6.55226 17.94 10.5579C17.94 14.5636 14.6383 17.8369 10.533 17.8369C6.42768 17.8369 3.12598 14.5636 3.12598 10.5579Z"></path>
				</svg>
				<h4>Search</h4>
			</div>
		</Block>
			<Block style={{flex: 1}}>
				<iframe src={sidebarUrl}
				        height="100%"
				        width="100%"
				        frameBorder="0"
				/>
			</Block>
		</div>
	);
};

export default Sidebar;
