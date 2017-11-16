import * as V from 'picodom';
import * as domlib from '../domlib';
import View from '../view';
import {getStatus as getSocketStatus} from '../socket';

export class MenuView extends View {
	constructor () {
		super();
		this.elStatus = document.createElement('div');
	}

	render () {
		const socketStatus = getSocketStatus();
		let statusClass = 'status ';
		if (socketStatus !== 1) {
			// eslint-disable-next-line no-magic-numbers
			if (socketStatus === 0 || socketStatus === 2) {
				statusClass += 'connecting';
			} else {
				statusClass += 'offline';
			}
		}
		V.patch(this.vStatus, this.vStatus = V.h('div', {'class': statusClass}), this.elStatus);

		if (!this.init) {
			domlib.setProps(this.el, {'class': 'ui fixed inverted menu'});
			const menuContainer = domlib.newAt(this.el, 'div', {'class': 'ui container'}),
				menuRight = domlib.newAt(menuContainer, 'div', {'class': 'menu right'});

			menuRight.appendChild(this.elStatus);
			this.init = true;
		}
	}
}
