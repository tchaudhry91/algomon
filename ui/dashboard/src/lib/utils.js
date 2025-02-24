export async function fetchDashboardData(fetcher) {
	const url = '/api/v1/checks';
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

export async function fetchNamedCheckData(fetcher, name) {
	const url = '/api/v1/checks/' + name;
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

export async function fetchNamedCheckFailures(fetcher, name) {
	const url = 'http://localhost:9967/api/v1/checks/' + name + '/failures';
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

export function getMinutesSinceDate(date) {
	const now = new Date();
	const diff = now - new Date(date);
	return Math.floor(diff / 60000);
}

export function getStatusIcon(status) {
	if (status === 'SUCCESSFUL') {
		return "<i class='fa-solid fa-check' style='color: #63E6BE;'></i>";
	} else {
		return "<i class='fa-solid fa-xmark' style='color: #df0c0c;'></i>";
	}
}
