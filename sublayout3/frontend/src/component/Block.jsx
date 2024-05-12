import React from 'react';
import style from './Block.module.scss';
import cls from '../util/cls.js';

const Block = (props) => {
	return (
		<div className={cls(style.block, props.class)} style={props.style}>
			{props.children}
		</div>
	);
};

export default Block;
