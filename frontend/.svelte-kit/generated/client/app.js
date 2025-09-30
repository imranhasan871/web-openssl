export { matchers } from './matchers.js';

export const nodes = [
	() => import('./nodes/0'),
	() => import('./nodes/1'),
	() => import('./nodes/2'),
	() => import('./nodes/3'),
	() => import('./nodes/4'),
	() => import('./nodes/5'),
	() => import('./nodes/6'),
	() => import('./nodes/7'),
	() => import('./nodes/8'),
	() => import('./nodes/9'),
	() => import('./nodes/10'),
	() => import('./nodes/11'),
	() => import('./nodes/12'),
	() => import('./nodes/13'),
	() => import('./nodes/14'),
	() => import('./nodes/15'),
	() => import('./nodes/16'),
	() => import('./nodes/17'),
	() => import('./nodes/18')
];

export const server_loads = [];

export const dictionary = {
		"/": [3],
		"/dashboard": [4,[2]],
		"/dashboard/certificates": [5,[2]],
		"/dashboard/certificates/convert": [6,[2]],
		"/dashboard/certificates/csr": [7,[2]],
		"/dashboard/certificates/generate": [8,[2]],
		"/dashboard/certificates/keys": [9,[2]],
		"/dashboard/certificates/parse": [10,[2]],
		"/dashboard/certificates/verify": [11,[2]],
		"/dashboard/encryption": [12,[2]],
		"/dashboard/encryption/asymmetric": [13,[2]],
		"/dashboard/encryption/hash": [14,[2]],
		"/dashboard/encryption/symmetric": [15,[2]],
		"/dashboard/operations": [16,[2]],
		"/login": [17],
		"/register": [18]
	};

export const hooks = {
	handleError: (({ error }) => { console.error(error) }),
	
	reroute: (() => {}),
	transport: {}
};

export const decoders = Object.fromEntries(Object.entries(hooks.transport).map(([k, v]) => [k, v.decode]));

export const hash = false;

export const decode = (type, value) => decoders[type](value);

export { default as root } from '../root.js';