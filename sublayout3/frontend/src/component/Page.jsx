import React from 'react';
import style from './Page.module.scss';
import NowPlaying from "./NowPlaying.jsx";

const Page = (props) => {
	return (
		<div className={style.page}>
			{props.children}
			<NowPlaying/>
		</div>
	);
};

export default Page;
