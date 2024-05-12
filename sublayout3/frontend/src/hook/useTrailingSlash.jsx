import React, { useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';

const useTrailingSlash = () => {
	const navigate = useNavigate();
	const location = useLocation();

	useEffect(() => {
		// Check if the current pathname ends with a trailing slash
		if (!location.pathname.endsWith('/')) {
			// Redirect to the same route with a trailing slash
			navigate(location.pathname + '/', {replace: true});
		}
	}, [history, location]);

	// Render null as empty hook
	return null;
};

export default useTrailingSlash;