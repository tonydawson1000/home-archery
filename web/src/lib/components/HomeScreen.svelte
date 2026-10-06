<script lang="ts">
	import type { SessionListRow } from '$lib/session/mappers';
	import type { PersonalBestCard } from '$lib/session/mappers';

	let {
		displayName,
		pb,
		sessions,
		onStart,
		onSwitch,
		onOpenSession,
		onHistory
	}: {
		displayName: string;
		pb: PersonalBestCard;
		sessions: SessionListRow[];
		onStart: () => void;
		onSwitch: () => void;
		onOpenSession: (id: string) => void;
		onHistory: () => void;
	} = $props();
</script>

<header>
	<h1>{displayName}</h1>
	<button class="hit" type="button" onclick={onSwitch}>Switch</button>
</header>
<section class="pb" aria-label="Personal best">
	<p>{pb.headline}</p>
	{#if pb.hint}
		<p class="hint">{pb.hint}</p>
	{/if}
</section>
<button class="hit primary" type="button" onclick={onStart}>Start session</button>
<ul>
	{#each sessions as row (row.id)}
		<li>
			<button class="hit" type="button" onclick={() => onOpenSession(row.id)}>
				{row.dateLabel}
				{row.arrowCount}
				{row.total}
				{#if row.isPersonalBest}★{/if}
			</button>
		</li>
	{/each}
</ul>
<button class="hit" type="button" onclick={onHistory}>History</button>
