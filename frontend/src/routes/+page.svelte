<script>
	import { onMount } from 'svelte';
	import { fetchTopics } from '$lib/api.js';
	import TopicCard from '$lib/components/TopicCard.svelte';
	import SearchBox from '$lib/components/SearchBox.svelte';

	let topics = [];
	let loading = true;
	let error = '';

	onMount(async () => {
		try {
			topics = await fetchTopics();
		} catch (e) {
			error = 'Could not load topics. Is the backend running?';
		}
		loading = false;
	});
</script>

<svelte:head>
	<title>ADHD Helper</title>
</svelte:head>

<div class="home">
	<header class="home-header fade-in">
		<h1 class="home-title">What would you like to explore?</h1>
		<p class="home-subtitle">Pick a topic to dive in, or search for something specific.</p>
	</header>

	<div class="search-section fade-in" style="animation-delay: 80ms">
		<SearchBox scope="" placeholder="Search across all topics..." />
	</div>

	{#if loading}
		<div class="loading-state">
			<div class="spinner"></div>
		</div>
	{:else if error}
		<div class="error-state fade-in">
			<p>{error}</p>
		</div>
	{:else}
		<div class="topic-grid">
			{#each topics as topic, i}
				<div class="fade-in" style="animation-delay: {120 + i * 60}ms">
					<TopicCard
						title={topic.title}
						description={topic.description}
						href="/{topic.slug}"
					/>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	.home {
		max-width: var(--max-width);
		margin: 0 auto;
		padding: var(--space-2xl) var(--space-lg) var(--space-3xl);
	}

	.home-header {
		text-align: center;
		margin-bottom: var(--space-xl);
	}

	.home-title {
		font-family: var(--font-heading);
		font-size: 2rem;
		font-weight: 700;
		color: var(--text-primary);
		margin-bottom: var(--space-sm);
		letter-spacing: -0.02em;
	}

	.home-subtitle {
		font-size: 1.05rem;
		color: var(--text-secondary);
	}

	.search-section {
		margin-bottom: var(--space-2xl);
	}

	.topic-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: var(--space-lg);
	}

	.loading-state {
		display: flex;
		justify-content: center;
		padding: var(--space-3xl);
	}

	.error-state {
		text-align: center;
		padding: var(--space-2xl);
		color: var(--text-secondary);
		background: var(--bg-card);
		border-radius: var(--radius-md);
		border: 1px solid var(--border-light);
	}

	@media (max-width: 640px) {
		.home {
			padding: var(--space-xl) var(--space-md) var(--space-2xl);
		}

		.home-title {
			font-size: 1.6rem;
		}

		.topic-grid {
			grid-template-columns: 1fr;
			gap: var(--space-md);
		}
	}
</style>
