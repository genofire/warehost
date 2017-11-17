/* eslint no-magic-numbers: "off"*/
/* eslint sort-keys: "off"*/

export default {
	'title': 'Warehost',
	// 'backend': 'wss://accounts.sum7.eu/ws'
	'backend': `ws${location.protocol === 'https:' ? 's' : ''}://${location.host}/ws`
};
