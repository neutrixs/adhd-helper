const BASE = '';

/**
 * Fetch all top-level topics
 * @returns {Promise<Array<{slug: string, title: string, description: string, isDir: boolean}>>}
 */
export async function fetchTopics() {
    const res = await fetch(`${BASE}/api/topics`);
    if (!res.ok) throw new Error('Failed to fetch topics');
    return res.json();
}

/**
 * Fetch a single topic detail by path
 * @param {string} path
 * @returns {Promise<{title: string, description: string, contentHtml: string, children: Array, breadcrumbs: Array}>}
 */
export async function fetchTopic(path) {
    const res = await fetch(`${BASE}/api/topics/${path}`);
    if (!res.ok) throw new Error('Failed to fetch topic');
    return res.json();
}

/**
 * Search topics recursively
 * @param {string} query
 * @param {string} [scope='']
 * @returns {Promise<Array<{path: string, title: string, snippet: string}>>}
 */
export async function searchTopics(query, scope = '') {
    const params = new URLSearchParams({ q: query });
    if (scope) params.set('scope', scope);
    const res = await fetch(`${BASE}/api/search?${params}`);
    if (!res.ok) throw new Error('Failed to search');
    return res.json();
}
