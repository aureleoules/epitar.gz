<script lang="ts">
	import '../app.css';
	import filesize from 'filesize';
	import HeartIcon from '../assets/svg/heart.svg?raw';

	export let stats;
</script>

<script context="module">
	import {variables} from '$lib/var';

	export async function load({fetch}) {
		const stats = await fetch(`${variables.apiUrl}/stats`)
		.then(res => res.json())
		.then(d => {
			return d;
		});

		return {
			props: {
				stats
			}
		};
	}

</script>

<nav>
	<ul>
		<li><strong>epitar.gz</strong></li>
		<li><a href="/">Search</a></li>
		<li><a href="/modules">Modules</a></li>
		<li><a href="/about">About</a></li>
	</ul>
	<ul>
		{#if stats}
		<li>
			<small>
				{stats.total_files.toLocaleString()} documents
			</small>
		</li>
		<li>
			<small>
				{filesize(stats.total_size)}
			</small>
		</li>
		{/if}
	</ul>
</nav>

<main>
	<slot />
</main>

<footer>
	<div class="footer-content">
		<small>
			built with {@html HeartIcon} by <a target="_blank" href="https://github.com/aureleoules">@aureleoules</a>
		</small>
	</div>
</footer>

<style lang="scss">

	nav {
		padding: 0 1em;
	}
	
	footer {
		.footer-content {
			position: absolute;
			bottom: 0;
			right: 0;
			margin-bottom: 15px;
			margin-right: 15px;
		}
	}
</style>
