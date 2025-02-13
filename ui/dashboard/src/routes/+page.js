import { fetchDashboardData } from '$lib/utils';

/** @type {import('./$types').PageLoad} */
export async function load({ params, fetch }) {
	let checks = await fetchDashboardData(fetch);
	return {
		checks: checks
	};
}
