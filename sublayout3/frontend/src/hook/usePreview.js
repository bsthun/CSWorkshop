import React, { useEffect, useMemo } from 'react';
import axios from "axios";
import { useParams } from "react-router-dom";

let state;
let lastFetch = null;

const getState = async () => new Promise((resolve, reject) => {
	if (new Date() - lastFetch > 5000) {
		axios('http://preview.prehack.sjpc.me:9001/api/preview/state').then((response) => {
			state = response.data;
			resolve(state);
		}).catch((error) => {
			console.log(error);
		});
	} else {
		resolve(state);
	}
});

const usePreview = (componentName) => {
	const {side} = useParams();
	const [state, setState] = React.useState(null);

	useEffect(() => {
		getState().then((state) => {
			setState(state);
		})
	}, []);

	return useMemo(() => {
		if (state === null) return;
		const sideInfo = state.data["side" + side.toUpperCase()];
		let url = sideInfo[componentName].url;
		const queryInfo = sideInfo[componentName].queryInfo;
		if (queryInfo) {
			// Random number between queryInfo.valueStart to queryInfo.valueEnd
			const value = Math.floor(Math.random() * (queryInfo.valueEnd - queryInfo.valueStart + 1)) + queryInfo.valueStart;
			url += "?" + queryInfo.key + "=" + value;
		}
		return url;
	}, [state]);
}

export default usePreview;
