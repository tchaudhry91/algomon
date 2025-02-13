<script>
	let { checks } = $props();
</script>

<div class="mx-auto is-size-2 mt-5 mb-2 has-text-centered is-uppercase">Current Checks</div>

<table class="table mx-auto">
	<thead>
		<tr>
			<td>Name</td>
			<td>Last Checked</td>
			<td>Status</td>
			<td>RC</td>
			<td>Actions</td>
		</tr>
	</thead>
	<tbody>
		{#each checks as check}
			<tr>
				<td>{check.name}</td>
				<td>{check.timestamp}</td>
				<td>{check.status}</td>
				<td>{check.rc}</td>
				<td>{check.action_keys}</td>
			</tr>
		{/each}
	</tbody>
</table>
