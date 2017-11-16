import * as V from 'picodom';
import View from '../view';


const DELAY_OF_NOTIFY = 15000,
	MAX_MESSAGE_SHOW = 10;

class NotifyView extends View {
	constructor () {
		super();
		if ('Notification' in window) {
			window.Notification.requestPermission();
		}

		this.messages = [];
		window.setInterval(this.removeLast.bind(this), DELAY_OF_NOTIFY);
	}

	removeLast () {
		this.messages.splice(0, 1);
		this.render();
	}

	renderMsg (msg) {
		const {messages} = this,
			content = [msg.content];

		let {render} = this;
		render = render.bind(this);

		if (msg.header) {
			content.unshift(V.h('div', {'class': 'header'}, msg.header));
		}


		return V.h(
			'div', {
				'class': `ui floating message ${msg.type}`
			},
			V.h('i', {
				'class': 'close icon',
				'onclick': () => {
					const index = messages.indexOf(msg);
					if (index !== -1) {
						messages.splice(index, 1);
						render();
					}
				}
			}), V.h('div', {'class': 'content'}, content)
		);
	}

	send (props, content) {
		let msg = props;
		if (typeof props === 'object') {
			msg.content = content;
		} else {
			msg = {
				'content': content,
				'type': props
			};
		}
		if ('Notification' in window &&
			window.Notification.permission === 'granted') {
			let body = msg.type,
				title = content;
			if (msg.header) {
				title = msg.header;
				body = msg.content;
			}

			// eslint-disable-next-line no-new
			new window.Notification(title, {
				'body': body,
				'icon': '/img/logo.jpg'
			});

			return;
		}
		if (this.messages.length > MAX_MESSAGE_SHOW) {
			this.removeLast();
		}

		this.messages.push(msg);
		this.render();
	}

	render () {
		V.patch(this.vel, this.vel = V.h('div', {'class': 'notifications'}, this.messages.map(this.renderMsg.bind(this))), this.el);
	}
}
// eslint-disable-next-line one-var
const singelton = new NotifyView();
export {singelton, NotifyView};
