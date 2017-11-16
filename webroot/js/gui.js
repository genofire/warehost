import * as domlib from './domlib';
import {MenuView} from './element/menu';
import Navigo from '../node_modules/navigo/lib/navigo';
import View from './view';
import {singelton as notify} from './element/notify';

const router = new Navigo(null, true, '#'),
	elMain = domlib.newEl('div', {'class': 'ui main container'}),
	elMenu = new MenuView();

export {router};

let init = false,
	currentView = new View();


export function render () {
	if (!document.body) {
		return;
	}

	if (!init) {
		notify.bind(document.body);
		elMenu.bind(document.body);

		document.body.appendChild(elMain);

		init = true;
	}

	notify.render();
	elMenu.render();

	router.resolve();
}


export function setView (toView) {
	currentView.unbind();
	currentView = toView;
	currentView.bind(elMain);
	currentView.render();
}
