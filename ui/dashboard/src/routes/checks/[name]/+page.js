import { fetchNamedCheckData, fetchNamedCheckFailures } from '$lib/utils';

/** @type {import('./$types').PageLoad} */
export async function load({ params, fetch }) {
	let check = await fetchNamedCheckData(fetch, params.name);
	let failures = await fetchNamedCheckFailures(fetch, params.name);
	return {
		check: check,
		failures: failures,
		name: params.name
	};
}
