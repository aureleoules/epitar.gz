<script lang="ts">
	import '../app.css';
	import filesize from 'filesize';
	import HeartIcon from '../assets/svg/heart.svg?raw';
	import { onMount } from 'svelte';

	import { variables } from '$lib/var';
	export let stats;

	onMount(async () => {
		stats = await fetch(`${variables.apiUrl}/stats`)
			.then((res) => res.json())
			.then((d) => {
				return d;
			})
			.catch((err) => {
				console.error(err);
				return {};
			});
	});
</script>

<nav>
	<ul>
		<li><strong>epitar.gz</strong></li>
		<li><a href="/">Documents</a></li>
		<li><a href="/news">News</a></li>
	</ul>
	<ul>
		<li><a href="/modules">Sources</a></li>
		{#if stats}
			<li>
				<small>
					{stats?.total_files?.toLocaleString()} documents
				</small>
			</li>
			<li>
				<small>
					{stats?.total_news?.toLocaleString()} news
				</small>
			</li>
			<li>
				<small>
					{stats?.total_size && filesize(stats.total_size)}
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
			built with {@html HeartIcon} by
			<a target="_blank" href="https://github.com/aureleoules">@aureleoules</a>
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
