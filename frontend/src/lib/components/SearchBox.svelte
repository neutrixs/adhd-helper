<script>
	import { searchTopics } from '$lib/api.js';
	import { goto } from '$app/navigation';

	export let scope = '';
	export let placeholder = 'Search topics...';

	let query = '';
	let results = [];
	let loading = false;
	let showDropdown = false;
	let debounceTimer;
	let inputEl;

	function handleInput() {
		clearTimeout(debounceTimer);
		if (!query.trim()) {
			results = [];
			showDropdown = false;
			loading = false;
			return;
		}
		loading = true;
		showDropdown = true;
		debounceTimer = setTimeout(async () => {
			try {
				results = await searchTopics(query, scope);
			} catch (e) {
				results = [];
			}
			loading = false;
		}, 300);
	}

	function handleBlur() {
		// Delay to allow click on results
		setTimeout(() => {
			showDropdown = false;
		}, 200);
	}

	function handleFocus() {
		if (query.trim() && (results.length > 0 || loading)) {
			showDropdown = true;
		}
	}

	function navigateResult(path) {
		query = '';
		results = [];
		showDropdown = false;
		goto(path);
	}

	function handleKeydown(e) {
		if (e.key === 'Escape') {
			showDropdown = false;
			inputEl?.blur();
		}
	}
</script>

<div class="search-wrapper" role="search">
	<div class="search-input-wrap" class:focused={showDropdown}>
		<svg class="search-icon" width="18" height="18" viewBox="0 0 18 18" fill="none">
			<circle cx="7.5" cy="7.5" r="5.5" stroke="currentColor" stroke-width="1.5"/>
			<path d="M12 12l4 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
		</svg>
		<input
			bind:this={inputEl}
			type="text"
			bind:value={query}
			on:input={handleInput}
			on:blur={handleBlur}
			on:focus={handleFocus}
			on:keydown={handleKeydown}
			{placeholder}
			aria-label="Search"
			autocomplete="off"
		/>
		{#if loading}
			<div class="search-spinner">
				<div class="spinner"></div>
			</div>
		{/if}
	</div>

	{#if showDropdown}
		<div class="search-dropdown fade-in">
			{#if loading}
				<div class="search-status">Searching...</div>
			{:else if results.length === 0 && query.trim()}
				<div class="search-status">No results found</div>
			{:else}
				{#each results as result}
					<button
						class="search-result"
						on:click={() => navigateResult(result.path)}
						type="button"
					>
						<span class="result-title">{result.title}</span>
						{#if result.snippet}
							<span class="result-snippet">{result.snippet}</span>
						{/if}
					</button>
				{/each}
			{/if}
		</div>
	{/if}
</div>

<style>
	.search-wrapper {
		position: relative;
		width: 100%;
		max-width: 560px;
		margin: 0 auto;
	}

	.search-input-wrap {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		background: var(--bg-card);
		border: 1.5px solid var(--border);
		border-radius: var(--radius-full);
		padding: var(--space-sm) var(--space-md);
		transition: border-color var(--duration-fast) var(--ease-out),
			box-shadow var(--duration-fast) var(--ease-out);
	}

	.search-input-wrap.focused,
	.search-input-wrap:focus-within {
		border-color: var(--accent);
		box-shadow: 0 0 0 3px var(--accent-light);
	}

	.search-icon {
		color: var(--text-muted);
		flex-shrink: 0;
	}

	input {
		flex: 1;
		border: none;
		background: none;
		font-family: var(--font-body);
		font-size: 0.95rem;
		color: var(--text-primary);
		outline: none;
		min-width: 0;
	}

	input::placeholder {
		color: var(--text-muted);
	}

	.search-spinner {
		flex-shrink: 0;
	}

	.search-dropdown {
		position: absolute;
		top: calc(100% + var(--space-sm));
		left: 0;
		right: 0;
		background: var(--bg-card);
		border: 1px solid var(--border);
		border-radius: var(--radius-md);
		box-shadow: var(--shadow-lg);
		max-height: 320px;
		overflow-y: auto;
		z-index: 100;
	}

	.search-status {
		padding: var(--space-md) var(--space-lg);
		color: var(--text-muted);
		font-size: 0.9rem;
		text-align: center;
	}

	.search-result {
		display: block;
		width: 100%;
		text-align: left;
		background: none;
		border: none;
		border-bottom: 1px solid var(--border-light);
		padding: var(--space-md) var(--space-lg);
		font-family: var(--font-body);
		cursor: pointer;
		transition: background var(--duration-fast) var(--ease-out);
	}

	.search-result:last-child {
		border-bottom: none;
	}

	.search-result:hover {
		background: var(--accent-light);
	}

	.result-title {
		display: block;
		font-family: var(--font-heading);
		font-weight: 600;
		font-size: 0.95rem;
		color: var(--text-primary);
		margin-bottom: 2px;
	}

	.result-snippet {
		display: block;
		font-size: 0.82rem;
		color: var(--text-secondary);
		line-height: 1.4;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
