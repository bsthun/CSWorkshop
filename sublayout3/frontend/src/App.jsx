import { BrowserRouter, Route, Router, Routes } from "react-router-dom";
import Preview from "./preview/index.jsx";
import Home from "./page/Home.jsx";
import "./style/index.scss"
import Search from "./page/Search.jsx";
import PreSearch from "./page/PreSearch.jsx";
import Artist from "./page/Artist.jsx";
import ArtistDiscography from "./page/ArtistDiscography.jsx";
import Generic from "./page/Generic.jsx";

function App() {
	return (
		<>
			<BrowserRouter>
				<Routes>
					<Route path="/side/:side/" caseSensitive={true}>
						<Route path="" element={<Preview/>}/>
						<Route path="home" element={<Home/>}/>
						<Route path="presearch" element={<PreSearch/>}/>
						<Route path="search" element={<Search/>}/>
						<Route path="artist" element={<Artist/>}/>
						<Route path="artistdiscography" element={<ArtistDiscography/>}/>
						<Route path="artistconcertlist" element={<Generic componentName="artistconcertlist" />}/>
						<Route path="concertdetail" element={<Generic componentName="concertdetail" />}/>
						<Route path="playlist" element={<Generic componentName="playlist" />}/>
						<Route path="track" element={<Generic componentName="track" />}/>
						<Route path="trackdetail" element={<Generic componentName="trackdetail" />}/>
						<Route path="podcasthome" element={<Generic componentName="podcasthome" />}/>
						<Route path="podcastgenre" element={<Generic componentName="podcastgenre" />}/>
						<Route path="podcastshow" element={<Generic componentName="podcastshow" />}/>
						<Route path="podcastepisode" element={<Generic componentName="podcastepisode" />}/>
						<Route path="album" element={<Generic componentName="album" />}/>
					</Route>
				</Routes>
			</BrowserRouter>
		</>
	)
}

export default App
