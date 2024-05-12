import { BrowserRouter, Route, Router, Routes } from "react-router-dom";
import Preview from "./preview/index.jsx";
import Home from "./page/Home.jsx";
import "./style/index.scss"
import Search from "./page/Search.jsx";
import PreSearch from "./page/PreSearch.jsx";
import Artist from "./page/Artist.jsx";
import ArtistDiscography from "./page/ArtistDiscography.jsx";

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
					</Route>
				</Routes>
			</BrowserRouter>
		</>
	)
}

export default App
