<script>
	import Dashboard from '$lib/Dashboard.svelte';
	let { data } = $props();
</script>

<div class="container">
	<Dashboard checks={data.checks} />
</div>
