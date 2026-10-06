<script lang="ts">
	import type { SessionListRow } from '$lib/session/mappers';

	let {
		sessions,
		onOpen,
		onHome
	}: {
		sessions: SessionListRow[];
		onOpen: (id: string) => void;
		onHome: () => void;
	} = $props();
</script>

<header>
	<h1>Sessions</h1>
	<button class="hit" type="button" onclick={onHome}>Home</button>
</header>
<ul>
	{#each sessions as row (row.id)}
		<li>
			<button class="hit" type="button" onclick={() => onOpen(row.id)}>
				{row.dateLabel}
				{row.arrowCount} arr
				{row.total}
				{#if row.isPersonalBest}★{/if}
			</button>
		</li>
	{/each}
</ul>
