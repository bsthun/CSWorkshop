import React from 'react';
import Sidebar from "../component/Sidebar.jsx";
import Page from "../component/Page.jsx";
import Block from "../component/Block.jsx";
import usePreview from "../hook/usePreview.js";

const Search = () => {
	const searchUrl = usePreview("search")
	return (
		<Page>
			<div style={{display: "flex", flexDirection: "row", gap: "0.5rem", padding: "0.5rem", height: "100%"}}>
				<Sidebar/>
				<div style={{
					flex: 7,
					display: "flex",
					flexDirection: "column",
					alignItems: "stretch",
					justifyContent: "start",
					gap: "0.5rem",
					height: "100%"
				}}>
					<Block style={{height: "100%"}}>
						<iframe src={searchUrl}
						        height="100%"
						        width="100%"
						        frameBorder="0"
						/>
					</Block>
				</div>
			</div>
		</Page>
	);
};

export default Search;
