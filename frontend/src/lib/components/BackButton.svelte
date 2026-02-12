<script>
	export let title = '';
	export let href = '/';
</script>

<nav class="top-bar">
	<a class="back-btn" {href}>
		<svg width="10" height="16" viewBox="0 0 10 16" fill="none">
			<path d="M8.5 1L1.5 8l7 7" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
		</svg>
		<span class="back-text">Back</span>
	</a>
	<h1 class="bar-title">{title}</h1>
	<div class="bar-spacer"></div>
</nav>

<style>
	.top-bar {
		display: flex;
		align-items: center;
		height: var(--nav-height);
		padding: 0 var(--space-md);
		background: var(--bg-primary);
		border-bottom: 1px solid var(--border-light);
		position: sticky;
		top: 0;
		z-index: 50;
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		background: rgba(253, 246, 236, 0.85);
	}

	.back-btn {
		display: flex;
		align-items: center;
		gap: var(--space-xs);
		color: var(--accent);
		text-decoration: none;
		font-family: var(--font-body);
		font-size: 0.95rem;
		font-weight: 400;
		padding: var(--space-sm) var(--space-sm) var(--space-sm) 0;
		flex-shrink: 0;
		min-width: 64px;
		transition: opacity var(--duration-fast) var(--ease-out);
	}

	.back-btn:active {
		opacity: 0.6;
	}

	.bar-title {
		flex: 1;
		text-align: center;
		font-family: var(--font-heading);
		font-size: 1.05rem;
		font-weight: 600;
		color: var(--text-primary);
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
		padding: 0 var(--space-sm);
	}

	.bar-spacer {
		min-width: 64px;
		flex-shrink: 0;
	}

	@media (min-width: 768px) {
		.top-bar {
			padding: 0 var(--space-lg);
		}
	}
</style>
