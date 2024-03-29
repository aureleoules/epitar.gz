<script lang="ts" context="module">
	import InfiniteScroll from 'svelte-infinite-scroll';
	import { variables } from '$lib/var';
	export async function load({ fetch }) {
		const modules = await fetch(`${variables.apiUrl}/modules`)
			.then((res) => res.json())
			.then((d) => {
				return d;
			})
			.catch((e) => {
				console.error(e);
				return [];
			});
		return {
			props: { modules }
		};
	}
</script>

<script lang="ts">
	import { onMount } from 'svelte';

	export let modules = [];
	let page = 1;
	let query = '';
	let module = '';
	let results = [];

	onMount(async () => {
		// Retrieve url params
		const urlParams = new URLSearchParams(window.location.search);
		query = urlParams.get('q') || "";
		module = urlParams.get('module') || "";
		
		// If query is set, search
		if (query) {
			search();
		}
	});

	function search(concat: boolean = false) {
		history.pushState(null, null, `?q=${query}&module=${module}`);

		fetch(`${variables.apiUrl}/documents/search?q=${query}&module=${module}&page=${page}`)
			.then((res) => {
				if (res.ok) {
					return res.json();
				}
				return [];
			})
			.then((res) => {
				if (concat) results = results.concat(res);
				else results = res;
				console.log(res);
			})
			.catch((err) => {
				console.log(err);
			});
	}

	function setParams(q: string, m: string) {
		query = q;
		module = m;
		page = 1;
		search();
	}

	function onSubmit(e: Event) {

	}
</script>

<svelte:head>
	<title>epitar.gz search index</title>
</svelte:head>

<div class="search-page">
	<article class="search-box">
		<header>Search for documents</header>
		<form on:submit={onSubmit}>
			<label for="q">Query</label>
			<input
				id="search"
				type="text"
				name="q"
				placeholder="thl, assembly, mathematics..."
				required
				value={query}
				on:input={(e) => setParams(e.target.value, module)}
			/>
			<label for="module">Source</label>
			<select value={module} on:input={(e) => setParams(query, e.target.value)} name="module">
				<option value="">All</option>
				{#each modules as module}
					<option value={module.slug}>{module.name}</option>
				{/each}
			</select>
			<button type="submit">Search</button>
		</form>
	</article>

	{#if results?.length}
		<div class="results grid">
			{#each results as result}
				<article>
					<header title={result.name}>
						<a target="_blank" href={result.origins[0].original_url}>{result.name}</a>
					</header>
					<p>
						{result.summary}
					</p>
					<footer>
						{#each result.origins as origin}
							<a role="button" target="_blank" href={origin.original_url}>{origin.module}</a>
						{/each}
					</footer>
				</article>
			{/each}

			<InfiniteScroll
				hasMore={true}
				window={true}
				threshold={100}
				on:loadMore={() => {
					page++;
					search(true);
				}}
			/>
		</div>
	{/if}
</div>

<style lang="scss">
	.search-page {
		padding: 1rem;
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

			footer {
				a {
					margin-right: 14px;
					padding: 8px;
					font-size: 16px;
				}
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
