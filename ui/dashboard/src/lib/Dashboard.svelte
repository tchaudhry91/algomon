<script>
	import { getMinutesSinceDate, getStatusIcon } from './utils';
	let { checks } = $props();
</script>

<div class="mx-auto is-size-1 mt-5 mb-5 has-text-centered">Dashboard</div>

<table class="table mx-auto">
	<thead>
		<tr>
			<td>Name</td>
			<td>Last Checked</td>
			<td>Status</td>
			<td>Actions</td>
		</tr>
	</thead>
	<tbody>
		{#each checks as check}
			<tr>
				<td><a href="/checks/?name={check.name}">{check.name}</a></td>
				<td align="center">{getMinutesSinceDate(check.timestamp)}m</td>
				<td align="center">{@html getStatusIcon(check.status)}</td>
				<td>{check.action_keys}</td>
			</tr>
		{/each}
	</tbody>
</table>
