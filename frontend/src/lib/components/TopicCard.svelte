<script>
	export let title = "";
	export let description = "";
	export let href = "/";
</script>

<a class="topic-card" {href}>
	<h3 class="card-title">{title}</h3>
	{#if description}
		<p class="card-desc">{description}</p>
	{/if}
	<div class="card-arrow">
		<svg width="20" height="20" viewBox="0 0 20 20" fill="none">
			<path
				d="M7 4l6 6-6 6"
				stroke="currentColor"
				stroke-width="1.5"
				stroke-linecap="round"
				stroke-linejoin="round"
			/>
		</svg>
	</div>
</a>

<style>
	.topic-card {
		display: block;
		position: relative;
		height: 100%;
		background: var(--bg-card);
		border: 1px solid var(--border-light);
		border-radius: var(--radius-md);
		padding: var(--space-lg);
		text-decoration: none;
		color: inherit;
		box-shadow: var(--shadow-sm);
		transition:
			transform var(--duration-normal) var(--ease-out),
			box-shadow var(--duration-normal) var(--ease-out),
			background var(--duration-normal) var(--ease-out),
			border-color var(--duration-normal) var(--ease-out);
		overflow: hidden;
	}

	.topic-card::before {
		content: "";
		position: absolute;
		top: 0;
		left: 0;
		width: 4px;
		height: 100%;
		background: var(--accent);
		opacity: 0;
		transition: opacity var(--duration-normal) var(--ease-out);
	}

	.topic-card:hover {
		transform: translateY(-2px);
		box-shadow: var(--shadow-hover);
		background: var(--bg-card-hover);
		border-color: var(--border);
	}

	.topic-card:hover::before {
		opacity: 1;
	}

	.topic-card:active {
		transform: translateY(0);
		box-shadow: var(--shadow-md);
	}

	.card-title {
		font-family: var(--font-heading);
		font-size: 1.15rem;
		font-weight: 600;
		margin-bottom: var(--space-sm);
		color: var(--text-primary);
		padding-right: var(--space-xl);
	}

	.card-desc {
		font-size: 0.92rem;
		color: var(--text-secondary);
		line-height: 1.55;
		display: -webkit-box;
		-webkit-line-clamp: 3;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}

	.card-arrow {
		position: absolute;
		top: var(--space-lg);
		right: var(--space-md);
		color: var(--text-muted);
		transition:
			color var(--duration-fast) var(--ease-out),
			transform var(--duration-normal) var(--ease-out);
	}

	.topic-card:hover .card-arrow {
		color: var(--accent);
		transform: translateX(3px);
	}
</style>
