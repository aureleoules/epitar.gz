<script context="module" lang="ts">
	export const prerender = true;
</script>

<script lang="ts">
	import {variables} from "$lib/var";

	let query = "";
	let results = [];

	function setSearch(q: string) {
		query = q;
		console.log(query);
		window.history.pushState({}, "", `?q=${q}`);

		// fetch data
		fetch(`${variables.apiUrl}/search?q=${q}`)
			.then(res => res.json())
			.then(data => {
				results = data;
				console.log(data);
			}).catch(err => {
				console.log(err);
			});
	}

</script>

<svelte:head>
	<title>epitar.gz search index</title>
</svelte:head>


<div>
	<input type="text" placeholder="Search" on:input={e => setSearch(e.target.value)} />
</div>

{#if results.length}
	<div class="results">
		{#each results as result}
			<a href={`${variables.apiUrl}/file/${result.id}`}>{result.name}</a>
		{/each}
	</div>
{/if}

<style lang="scss">
	.results {
		display: flex;
		flex-direction: column;
		margin-top: 1em;
	}	
</style>
