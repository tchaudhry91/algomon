<script>
	let { data } = $props();
	import CheckHistory from '$lib/CheckHistory.svelte';
</script>

<div class="container">
	<CheckHistory check={data.check} failures={data.failures} name={data.name} />
</div>
