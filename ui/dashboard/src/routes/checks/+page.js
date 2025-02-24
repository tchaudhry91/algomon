import { fetchNamedCheckData, fetchNamedCheckFailures } from '$lib/utils';

/** @type {import('./$types').PageLoad} */
export async function load({ params, fetch, url }) {
	let name = url.searchParams.get('name');
	let check = await fetchNamedCheckData(fetch, name);
	let failures = await fetchNamedCheckFailures(fetch, name);
	return {
		check: check,
		failures: failures,
		name: name.name
	};
}
