<script lang="ts">
	import { variables } from '$lib/var';
	import InfiniteScroll from 'svelte-infinite-scroll';
	import { onMount } from 'svelte';

	let newsgroups = [];
	let page = 1;
	let query = '';
	let newsgroup = '';
	let results = [];

	onMount(async () => {
		newsgroups = await fetch(`${variables.apiUrl}/newsgroups`)
			.then((res) => res.json())
			.then((d) => {
				return d;
			})
			.catch((e) => {
				console.error(e);
				return [];
			});

		// Retrieve url params
		const urlParams = new URLSearchParams(window.location.search);
		query = urlParams.get('q') || "";
		newsgroup = urlParams.get('newsgroup') || "";
		
		// If query is set, search
		if (query) {
			search();
		}
	});

	function search(concat: boolean = false) {
		history.pushState(null, null, `/news?q=${query}&newsgroup=${newsgroup}`);

		fetch(`${variables.apiUrl}/news/search?q=${query}&newsgroup=${newsgroup}&page=${page}`)
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
		newsgroup = m;
		page = 1;
		search();
	}

	function onSubmit(e: Event) {
		e.preventDefault();
	}
</script>

<svelte:head>
	<title>epitar.gz search index</title>
</svelte:head>

<div class="search-page">
	<article class="search-box">
		<header>Search for news</header>
		<form on:submit={onSubmit}>
			<label for="q">Query</label>
			<input
				type="text"
				name="q"
				placeholder="piscine, traces evalexpr, 42sh..."
				required
				value={query}
				on:input={(e) => setParams(e.target.value, newsgroup)}
			/>
			<label for="newsgroup">Newsgroup</label>
			<select value={newsgroup} on:input={(e) => setParams(query, e.target.value)} name="newsgroup">
				<option value="">All</option>
				{#each newsgroups as newsgroup}
					<option value={newsgroup}>{newsgroup}</option>
				{/each}
			</select>
			<button type="submit">Search</button>
		</form>
	</article>

	{#if results?.length}
		<div class="results grid">
			{#each results as result}
				<article>
					<header title={result.subject}>
						<a target="_blank" href={`https://news.infinity.study/?news=${result.message_id}`}>{result.subject}</a>
					</header>
					<p>
						{result.summary}
					</p>
					<footer>
						{#each result.newsgroups.split(",") as g}
							<a role="button" target="_blank" href={`https://news.infinity.study/?news=${result.message_id}`}>{g}</a>
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
					font-size: 13px;
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
