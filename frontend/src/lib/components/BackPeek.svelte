<script>
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { onMount, onDestroy } from 'svelte';

	export let backHref = '/';

	let swipeProgress = 0; // 0 to 1
	let isSwiping = false;
	let startX = 0;
	let startY = 0;
	let currentX = 0;
	let isTouch = false;
	let containerEl;
	let directionLocked = false;
	let isHorizontal = false;

	const EDGE_ZONE = 30; // px from left edge
	const THRESHOLD = 0.35;
	const SCREEN_WIDTH = typeof window !== 'undefined' ? window.innerWidth : 400;

	onMount(() => {
		isTouch = 'ontouchstart' in window;
	});

	function handleTouchStart(e) {
		if (!isTouch) return;
		const touch = e.touches[0];
		if (touch.clientX > EDGE_ZONE) return;

		startX = touch.clientX;
		startY = touch.clientY;
		isSwiping = true;
		directionLocked = false;
		isHorizontal = false;
		swipeProgress = 0;
		document.body.classList.add('swiping');
	}

	function handleTouchMove(e) {
		if (!isSwiping) return;
		const touch = e.touches[0];
		const deltaX = touch.clientX - startX;
		const deltaY = touch.clientY - startY;

		// Lock direction after first significant movement
		if (!directionLocked && (Math.abs(deltaX) > 10 || Math.abs(deltaY) > 10)) {
			directionLocked = true;
			isHorizontal = Math.abs(deltaX) > Math.abs(deltaY);
		}

		if (!isHorizontal) {
			isSwiping = false;
			swipeProgress = 0;
			return;
		}

		e.preventDefault();
		currentX = Math.max(0, deltaX);
		swipeProgress = Math.min(1, currentX / (SCREEN_WIDTH * 0.6));
	}

	function handleTouchEnd() {
		if (!isSwiping) return;
		isSwiping = false;
		document.body.classList.remove('swiping');

		if (swipeProgress >= THRESHOLD) {
			// Complete the swipe
			swipeProgress = 1;
			setTimeout(() => {
				goto(backHref);
				swipeProgress = 0;
			}, 250);
		} else {
			swipeProgress = 0;
		}
	}
</script>

<div
	class="back-peek-container"
	bind:this={containerEl}
	on:touchstart={handleTouchStart}
	on:touchmove={handleTouchMove}
	on:touchend={handleTouchEnd}
	role="presentation"
>
	<!-- Behind layer (previous page hint) -->
	{#if swipeProgress > 0}
		<div
			class="peek-behind"
			style="opacity: {swipeProgress * 0.6}"
		>
			<div class="peek-behind-content">
				<svg width="24" height="24" viewBox="0 0 24 24" fill="none" style="opacity: {Math.min(1, swipeProgress * 2)}">
					<path d="M15 18l-6-6 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
				</svg>
			</div>
		</div>
	{/if}

	<!-- Current page layer -->
	<div
		class="peek-current"
		class:swiping={isSwiping || swipeProgress > 0}
		style="transform: translateX({swipeProgress * 80}%); {swipeProgress > 0 ? '' : ''}"
	>
		<slot />
	</div>

	<!-- Edge shadow during swipe -->
	{#if swipeProgress > 0}
		<div
			class="peek-edge-shadow"
			style="opacity: {swipeProgress * 0.4}; left: calc({swipeProgress * 80}% - 20px)"
		></div>
	{/if}
</div>

<style>
	.back-peek-container {
		position: relative;
		min-height: 100dvh;
	}

	:global(body.swiping) {
		overflow-x: clip;
	}

	.peek-current {
		position: relative;
		min-height: 100dvh;
		background: var(--bg-primary);
		z-index: 2;
		will-change: transform;
	}

	.peek-current.swiping {
		transition: none;
	}

	.peek-current:not(.swiping) {
		transition: transform var(--duration-slow) var(--ease-out);
	}

	.peek-behind {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background: var(--bg-secondary);
		z-index: 1;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.peek-behind-content {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 48px;
		height: 48px;
		border-radius: 50%;
		background: var(--bg-card);
		box-shadow: var(--shadow-md);
		color: var(--accent);
	}

	.peek-edge-shadow {
		position: fixed;
		top: 0;
		width: 20px;
		height: 100%;
		background: linear-gradient(to right, rgba(59, 47, 47, 0.15), transparent);
		z-index: 3;
		pointer-events: none;
	}
</style>
