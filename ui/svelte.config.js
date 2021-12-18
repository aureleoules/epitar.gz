import adapter from '@sveltejs/adapter-node';
import preprocess from 'svelte-preprocess';
import replace from '@rollup/plugin-replace';
import {config} from 'dotenv';

const production = !process.env.ROLLUP_WATCH;

/** @type {import('@sveltejs/kit').Config} */
const svelteConfig = {
	// Consult https://github.com/sveltejs/svelte-preprocess
	// for more information about preprocessors
	preprocess: preprocess(),

	kit: {
		adapter: adapter(),

		// hydrate the <div id="svelte"> element in src/app.html
		target: '#svelte'
	},
	plugins: [
		replace({
			// stringify the object       
			__myapp: JSON.stringify({
				env: {
					isProd: production,
					...config().parsed // attached the .env config
				}
			}),
		}),
	],
};

export default svelteConfig;
