<script lang="ts">
	import { variables } from '$lib/var';
	import { onMount } from 'svelte';

	let modules = [];

	onMount(async () => {
		modules = await fetch(`${variables.apiUrl}/modules`)
			.then((res) => res.json())
			.then((d) => {
				return d;
			})
			.catch((e) => {
				console.error(e);
				return [];
			});
	});
</script>

<div class="container">
	<h1>Indexed sources</h1>
	<table>
		<thead>
			<tr>
				<th>Name</th>
				<th>Description</th>
				<th>Archived by</th>
			</tr>
		</thead>
		<tbody>
			{#each modules as module}
				<tr>
					<td><a target="_blank" href={module.url}>{module.name}</a></td>
					<!-- <td>{filesize(modulesize)}</td> -->
					<td>{module.description}</td>
					<td>{module.authors.map((x) => x.name).join(', ')}</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
