<script>
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { fetchTopic } from '$lib/api.js';
	import BackButton from '$lib/components/BackButton.svelte';
	import BackPeek from '$lib/components/BackPeek.svelte';
	import SearchBox from '$lib/components/SearchBox.svelte';
	import TopicCard from '$lib/components/TopicCard.svelte';

	let topic = null;
	let loading = true;
	let error = '';

	// Reactive: re-fetch when path changes
	$: topicPath = $page.params.path || '';
	$: {
		if (topicPath !== undefined) {
			loadTopic(topicPath);
		}
	}

	async function loadTopic(path) {
		loading = true;
		error = '';
		topic = null;
		try {
			topic = await fetchTopic(path);
		} catch (e) {
			error = 'Could not load this topic.';
		}
		loading = false;
	}

	function backHref(breadcrumbs) {
		if (!breadcrumbs || breadcrumbs.length <= 1) return '/';
		return breadcrumbs[breadcrumbs.length - 2].path;
	}
</script>

<svelte:head>
	{#if topic}
		<title>{topic.title} — ADHD Helper</title>
	{:else}
		<title>ADHD Helper</title>
	{/if}
</svelte:head>

<BackPeek backHref={topic ? backHref(topic.breadcrumbs) : '/'}>
	<BackButton
		title={topic ? topic.title : ''}
		href={topic ? backHref(topic.breadcrumbs) : '/'}
	/>

	<div class="topic-page">
		{#if loading}
			<div class="loading-state">
				<div class="spinner"></div>
			</div>
		{:else if error}
			<div class="error-state fade-in">
				<p>{error}</p>
			</div>
		{:else if topic}
			{#if topic.breadcrumbs && topic.breadcrumbs.length > 1}
				<nav class="breadcrumbs fade-in" aria-label="Breadcrumbs">
					<a href="/">Home</a>
					{#each topic.breadcrumbs as crumb, i}
						<span class="crumb-sep">/</span>
						{#if i < topic.breadcrumbs.length - 1}
							<a href={crumb.path}>{crumb.title}</a>
						{:else}
							<span class="crumb-current">{crumb.title}</span>
						{/if}
					{/each}
				</nav>
			{/if}

			<div class="search-section fade-in" style="animation-delay: 60ms">
				<SearchBox scope={topicPath} placeholder="Search within this topic..." />
			</div>

			{#if topic.contentHtml}
				<article class="markdown-body fade-in" style="animation-delay: 100ms">
					{@html topic.contentHtml}
				</article>
			{/if}

			{#if topic.children && topic.children.length > 0}
				<section class="children-section fade-in" style="animation-delay: 140ms">
					<h2 class="children-heading">Subtopics</h2>
					<div class="children-grid">
						{#each topic.children as child, i}
							<div class="fade-in" style="animation-delay: {180 + i * 50}ms">
								<TopicCard
									title={child.title}
									description={child.description}
									href="/{topicPath}/{child.slug}"
								/>
							</div>
						{/each}
					</div>
				</section>
			{/if}
		{/if}
	</div>
</BackPeek>

<style>
	.topic-page {
		max-width: var(--max-width);
		margin: 0 auto;
		padding: var(--space-lg) var(--space-lg) var(--space-3xl);
	}

	.breadcrumbs {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: var(--space-xs);
		font-size: 0.85rem;
		color: var(--text-muted);
		margin-bottom: var(--space-lg);
	}

	.breadcrumbs a {
		color: var(--accent);
		text-decoration: none;
	}

	.breadcrumbs a:hover {
		text-decoration: underline;
	}

	.crumb-sep {
		color: var(--border);
		margin: 0 2px;
	}

	.crumb-current {
		color: var(--text-secondary);
	}

	.search-section {
		margin-bottom: var(--space-xl);
		position: relative;
		z-index: 10;
	}

	.markdown-body {
		position: relative;
		z-index: 1;
		margin-bottom: var(--space-2xl);
	}

	.children-section {
		margin-top: var(--space-lg);
	}

	.children-heading {
		font-family: var(--font-heading);
		font-size: 1.2rem;
		font-weight: 600;
		color: var(--text-primary);
		margin-bottom: var(--space-md);
	}

	.children-grid {
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
		.topic-page {
			padding: var(--space-md) var(--space-md) var(--space-2xl);
		}

		.children-grid {
			grid-template-columns: 1fr;
			gap: var(--space-md);
		}
	}
</style>
