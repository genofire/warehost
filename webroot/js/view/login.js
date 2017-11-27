import * as domlib from '../domlib';
import * as gui from '../gui';
import * as socket from '../socket';
import * as store from '../store';
import View from '../view';

class LoginView extends View {

	login () {
		const username = this.inputUsername.value,
			password = this.inputPassword.value;
		this.btn.classList.add('loading');
		socket.sendjson({
			'body': {
				'password': password,
				'username': username
			},
			'subject': 'login'
		}, (msg) => {
			store.isLogin = msg.body;
			this.inputPassword.value = '';
			if (msg.subject === 'login' && msg.body) {
				socket.sendjson({'subject': 'auth_status'});
				this.form.classList.remove('top', 'attached');
				this.notify.classList.add('hidden');
			} else {
				this.form.classList.add('top', 'attached');
				this.notify.classList.remove('hidden');
				this.inputPassword.focus();
			}
			this.btn.classList.remove('loading');
		});
	}

	// eslint-disable-next-line class-methods-use-this
	render () {
		if (!this.init) {
			domlib.newAt(this.el, 'h2', {'class': 'ui header'}, 'Login');
			const event = this.login.bind(this),
				form = domlib.newAt(this.el, 'form', {
					'class': 'ui segment form',
					'onsubmit': event
				}),
				fieldUsername = domlib.newAt(form, 'div', {'class': 'field'}),
				fieldPassword = domlib.newAt(form, 'div', {'class': 'field'}),
				fieldInputUsername = domlib.newAt(fieldUsername, 'div', {'class': 'ui left icon input'}),
				fieldInputPassword = domlib.newAt(fieldPassword, 'div', {'class': 'ui left icon input'});


			domlib.newAt(fieldInputUsername, 'i', {'class': 'user icon'});
			this.inputUsername = domlib.newAt(fieldInputUsername, 'input', {
				'placeholder': 'Benutzername',
				'type': 'text'
			});

			domlib.newAt(fieldInputPassword, 'i', {'class': 'lock icon'});
			this.inputPassword = domlib.newAt(fieldInputPassword, 'input', {
				'placeholder': 'Passwort',
				'type': 'password'
			});
			this.btn = domlib.newAt(form, 'div', {
				'class': 'ui fluid large primary submit button',
				'onclick': event
			}, 'Login');

			this.notify = domlib.newAt(this.el, 'div', {'class': 'ui hidden bottom attached error message'}, 'Anmeldung fehlgeschlagen!');

			this.logout = domlib.newAt(this.el, 'div', {'class': 'ui info message'});

			this.form = form;
			this.init = true;
		}
		if (store.isLogin) {
			domlib.removeChild(this.form);
			domlib.removeChild(this.notify);
			domlib.appendChild(this.el, this.logout);
			this.logout.innerHTML = `Sie sind bereits eingeloggt, ${store.login.username}!`;
		} else {
			domlib.removeChild(this.logout);
			domlib.appendChild(this.el, this.form);
			domlib.appendChild(this.el, this.notify);
		}
	}
}

const loginView = new LoginView();

gui.router.on('/login', () => {
	gui.setView(loginView);
});
