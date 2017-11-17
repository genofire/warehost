import * as V from 'picodom';
import * as domlib from '../domlib';
import * as socket from '../socket';
import * as store from '../store';
import View from '../view';
import {singelton as notify} from './notify';
import {render} from '../gui';

export class MenuView extends View {
	constructor () {
		super();
		this.elStatus = document.createElement('div');
	}

	render () {
		const socketStatus = socket.getStatus();
		let statusClass = 'status ',
			vLogin = V.h('a', {
				'class': 'item',
				'href': '#/login'
			}, 'Login');
		if (socketStatus !== 1) {
			// eslint-disable-next-line no-magic-numbers
			if (socketStatus === 0 || socketStatus === 2) {
				statusClass += 'connecting';
			} else {
				statusClass += 'offline';
			}
		}
		if (store.isLogin) {
			vLogin = V.h('a', {
				'class': 'item',
				'href': '#/',
				'onclick': () => socket.sendjson({'subject': 'logout'}, (msg) => {
					if (msg.body) {
						store.isLogin = false;
						store.login = {};
						render();
					} else {
						notify.send({
							'header': 'Abmeldung ist fehlgeschlagen',
							'type': 'error'
						}, 'Logout');
					}
				})
			}, 'Logout');
		}

		V.patch(this.vStatus, this.vStatus = V.h('div', {'class': statusClass}), this.elStatus);


		if (!this.init) {
			domlib.setProps(this.el, {'class': 'ui fixed inverted menu'});
			const menuContainer = domlib.newAt(this.el, 'div', {'class': 'ui container'});
			this.menuRight = domlib.newAt(menuContainer, 'div', {'class': 'menu right'});

			this.menuRight.appendChild(this.elStatus);
			this.init = true;
		}

		V.patch(this.vLogin, this.vLogin = vLogin, this.menuRight);
	}
}
