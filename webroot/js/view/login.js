import * as dom from '../domlib';
import * as gui from '../gui';
import View from '../view';

class LoginView extends View {
	// eslint-disable-next-line class-methods-use-this
	render () {
		const h1 = dom.newAt(this.el, 'h1');
		h1.innerHTML = 'Login';
	}
}

const login = new LoginView();

gui.router.on('/login', () => {
	gui.setView(login);
});
