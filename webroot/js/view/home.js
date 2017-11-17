import * as dom from '../domlib';
import * as gui from '../gui';
import View from '../view';

class HomeView extends View {
	// eslint-disable-next-line class-methods-use-this
	render () {
		if (!this.init) {
			const h1 = dom.newAt(this.el, 'h1');
			h1.innerHTML = 'Home';
			this.init = true;
		}
	}
}

const homeView = new HomeView();

gui.router.on('/', () => {
	gui.setView(homeView);
});
