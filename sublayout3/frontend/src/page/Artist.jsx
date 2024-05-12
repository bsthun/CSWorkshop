import React from 'react';
import Sidebar from "../component/Sidebar.jsx";
import Page from "../component/Page.jsx";
import Block from "../component/Block.jsx";
import usePreview from "../hook/usePreview.js";

const Artist = () => {
	const aboutArtistUrl = usePreview("aboutartist")
	const artistConcertSectionUrl = usePreview("artistconcertsection")

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
					<Block style={{flex:1}}>
						<iframe src={aboutArtistUrl}
						        width="100%"
						        style={{flex: 1}}
						        frameBorder="0"
						/>
						<iframe src={artistConcertSectionUrl}
						        style={{height: 300, marginTop: 12}}
						        width="100%"
						        frameBorder="0"
						/>
					</Block>
				</div>
			</div>
		</Page>
	);
};

export default Artist;
