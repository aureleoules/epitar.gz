<script context="module" lang="ts">
	export const prerender = true;
</script>

<script lang="ts">
	import { variables } from '$lib/var';
	import { onMount } from 'svelte';

	let query = '';
	let results = [];

	onMount(() => {
		// get url param
		query = new URLSearchParams(window.location.search).get('q');
		if (query) {
			console.log('query', query);
			search(query);
		}
	});

	function search(q: string) {
		fetch(`${variables.apiUrl}/search?q=${query}`)
			.then((res) => {
				if (res.ok) {
					return res.json();
				}
				return [];
			})
			.then((res) => {
				results = res;
				console.log(res);
			})
			.catch((err) => {
				console.log(err);
			});
	}

	function setSearch(q: string) {
		query = q;
		window.history.pushState({}, '', `?q=${q}`);
		search(q);
	}
</script>

<svelte:head>
	<title>epitar.gz search index</title>
</svelte:head>

<div class="search-page">
	<article class="search-box">
		<header>epitar.gz search index</header>
		<form>
			<input
				type="text"
				name="q"
				placeholder="thl, assembly, mathematics..."
				required
				on:input={(e) => setSearch(e.target.value)}
			/>
			<button type="submit">Search</button>
		</form>
	</article>

	{#if results?.length}
		<div class="results grid">
			{#each results as result}
				<article>
					<header title={result.name}>
						<a target="_blankl" href={`${variables.apiUrl}/file/${result.id}`}>{result.name}</a>
					</header>
					<p>
						{result.summary}
					</p>
					<footer>
						<a target="_blankl" href={`${variables.apiUrl}/file/${result.id}`}>Download</a>
					</footer>
				</article>
			{/each}
		</div>
	{/if}
</div>

<style lang="scss">
	.search-page {
		padding: 1rem;
		min-height: 80vh;
		// display: flex;
		flex-direction: column;
		align-items: center;
	}

	.results {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;

		article {
			width: 450px;

			header {
				overflow: hidden;
				text-overflow: ellipsis;
				white-space: nowrap;
			}

			p {
				white-space: wrap;
				height: 100px;
			}
		}
	}

	article.search-box {
		width: 50%;
		margin: 25px auto;

		@media screen and (max-width: 1000px) {
			width: 100%;
		}
	}
</style>
