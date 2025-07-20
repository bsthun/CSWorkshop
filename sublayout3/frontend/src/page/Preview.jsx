import React, { useEffect } from 'react';
import Home from "./Home.jsx";
import PreSearch from "./PreSearch.jsx";
import Search from "./Search.jsx";
import Artist from "./Artist.jsx";
import ArtistDiscography from "./ArtistDiscography.jsx";
import Generic from "./Generic.jsx";

const Preview = () => {
	const comps = [
		() => <Home/>,
		() => <PreSearch/>,
		() => <Search/>,
		() => <Artist/>,
		() => <ArtistDiscography/>,
		() => <Generic componentName="artistconcertlist"/>,
		() => <Generic componentName="concertdetail"/>,
		() => <Generic componentName="playlist"/>,
		() => <Generic componentName="track"/>,
		() => <Generic componentName="trackdetail"/>,
		() => <Generic componentName="podcasthome"/>,
		() => <Generic componentName="podcastgenre"/>,
		() => <Generic componentName="podcastshow"/>,
		() => <Generic componentName="podcastepisode"/>,
		() => <Generic componentName="album"/>
	]

	const [index, useIndex] = React.useState(0)

	useEffect(() => {
		const interval = setInterval(() => {
			useIndex(index + 1)
		}, 5000)
		return () => clearInterval(interval)
	}, [])

	return (
		<div>
			{index}
			{comps[index % comps.length]()}
		</div>
	);
};

export default Preview;
