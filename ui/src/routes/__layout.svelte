<script lang="ts">
	import '../app.css';
</script>


<nav>
	<ul>
		<li><strong>epitar.gz</strong></li>
	</ul>
	<ul>
		<li><a href="/">Search</a></li>
		<li><a href="/about">About</a></li>
	</ul>
</nav>

<main>
	<slot />
</main>

<style>
	nav {
		padding: 0 1em;
	}
</style>
