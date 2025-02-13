export async function fetchDashboardData(fetcher) {
	const url = 'http://localhost:9967/api/v1/checks';
	try {
		const response = await fetcher(url);
		if (!response.ok) {
			throw new Error(`Response status: ${response.status}`);
		}
		const json = await response.json();
		return json;
	} catch (error) {
		console.error(error.message);
	}
}
