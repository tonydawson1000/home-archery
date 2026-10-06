<script lang="ts">
	import type { Session } from '$lib/api/types';
	import { endRow } from '$lib/session/mappers';

	let { session, onResume }: { session: Session; onResume: () => void } = $props();

	const started = $derived(
		new Intl.DateTimeFormat('en-GB', { day: 'numeric', month: 'short', year: 'numeric' }).format(
			new Date(session.startedAt)
		)
	);
	const { summary } = $derived(session);
</script>

<h1>{started}</h1>
<p>
	Total {summary.total} · {summary.arrowCount} arr · H {summary.hits} G {summary.golds}
</p>
{#each session.ends as end (end.endNumber)}
	{@const row = endRow(end)}
	<p>
		End {row.endNumber}
		{#each row.codes as code, i (i)}
			{code}
		{/each}
		{row.total}
	</p>
{/each}
{#if session.status === 'in_progress'}
	<button class="hit primary" type="button" onclick={onResume}>Resume</button>
{/if}
