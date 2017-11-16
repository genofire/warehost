export function setProps (el, props) {
	if (props) {
		if (props.class) {
			let classList = props.class;
			if (typeof props.class === 'string') {
				classList = classList.split(' ');
			}
			el.classList.add(...classList);
		}
	}
}

export function newEl (eltype, props, content) {
	const el = document.createElement(eltype);
	setProps(el, props);
	if (content) {
		el.innerHTML = content;
	}

	return el;
}

// eslint-disable-next-line max-params
export function newAt (at, eltype, props, content) {
	const el = document.createElement(eltype);
	setProps(el, props);
	if (content) {
		el.innerHTML = content;
	}
	at.appendChild(el);

	return el;
}
