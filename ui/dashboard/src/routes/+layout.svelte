<script>
	import Navbar from '$lib/Navbar.svelte';
	import '@fortawesome/fontawesome-free/css/all.min.css';

	let { children } = $props();
</script>

<Navbar />

{@render children()}
