import React from 'react';
import Block from "../component/Block.jsx";
import Page from "../component/Page.jsx";
import Sidebar from "../component/Sidebar.jsx";

const Home = () => {
	return (
		<Page>
			<div style={{display: "flex", flexDirection: "row", gap: "0.5rem", padding: "0.5rem", height: "100%"}}>
				<Sidebar />
				<div style={{
					flex: 5,
					display: "flex",
					flexDirection: "column",
					alignItems: "stretch",
					justifyContent: "start",
					gap: "0.5rem",
					height: "100%"
				}}>
					<Block>
						<iframe src="http://csc105-workshop.bsthun.com/components/comp01"
						        height="100%"
						        width="100%"
						        frameBorder="0"
						/>
					</Block>
					<Block style={{flex: 1}}>
						<iframe src="http://csc105-workshop.bsthun.com/components/comp01"
						        height="100%"
						        width="100%"
						        frameBorder="0"
						/>
					</Block>
				</div>
				<div style={{
					flex: 2,
					display: "flex",
					flexDirection: "column",
					alignItems: "stretch",
					justifyContent: "start",
					gap: "1rem",
				}}>
					<Block style={{
						height: "100%"
					}}>
						<iframe src="http://csc105-workshop.bsthun.com/components/comp01"
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

export default Home;
