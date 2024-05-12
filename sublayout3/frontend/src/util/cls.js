const cls = (...classes) => {
	  return classes.filter(Boolean).join(' ');
}

export default cls;