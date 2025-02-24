<script>
	let { name, check, failures } = $props();
	import { getStatusIcon, getMinutesSinceDate } from '$lib/utils';

	let showOutput = $state(false);
	let output = $state(check[0].combined_out);

	function showOutputModal(out) {
		let outP = out;
		try {
			outP = JSON.stringify(JSON.parse(out), null, 2);
		} catch (e) {
			console.log(e);
			outP = out;
		}
		output = outP;
		showOutput = true;
	}
</script>

<div class="mx-auto is-size-1 mt-5 mb-5 has-text-centered">{name}</div>
<div class="columns">
	<div class="column is-one-quarter">
		<div class="card mt-5">
			<div class="card-content">
				<p class="title">
					{check[0].name}
				</p>
				<p class="subtitle mt-5">
					{@html getStatusIcon(check[0].status)}
					Checked {getMinutesSinceDate(check[0].timestamp)}m ago
				</p>
			</div>
			<footer class="card-footer">
				<p class="card-footer-item">
					<button class="button is-light" onclick={() => showOutputModal(check[0].combined_out)}
						>Output</button
					>
				</p>
			</footer>
		</div>
	</div>
	<div class="column">
		<div class="mx-auto is-size-3 mt-5 mb-5 has-text-centered">Failure History</div>
		<table class="table mx-auto">
			<thead>
				<tr>
					<th align="center">Time</th>
					<th align="center">Status</th>
					<th align="center">Output</th>
					<th align="center">Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each failures as check}
					<tr>
						<td align="center">{getMinutesSinceDate(check.timestamp)}m ago</td>
						<td align="center">{@html getStatusIcon(check.status)}</td>
						<td align="center">
							<button
								class="button"
								onclick={() => {
									showOutputModal(check.combined_out);
								}}>Output</button
							>
						</td>
						<td>{check.action_keys}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
	<div class="column is-one-quarter"></div>
</div>

<div class="modal {showOutput ? 'is-active' : ''}">
	<div class="modal-background"></div>
	<div class="modal-content">
		<pre>
			<code>
				{output}
			</code>
		</pre>
	</div>
	<button class="modal-close is-large" aria-label="close" onclick={() => (showOutput = false)}
	></button>
</div>
