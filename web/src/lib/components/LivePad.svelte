<script lang="ts">
	import type { LivePadView } from '$lib/session/mappers';
	import type { ScoreCode } from '$lib/api/types';
	import { summaryLine } from '$lib/session/mappers';

	const pad: ScoreCode[][] = [
		['X', '10', '9'],
		['8', '7', '6'],
		['5', '4', '3'],
		['2', '1', 'M']
	];

	let {
		view,
		canFinish,
		onScore,
		onFinish
	}: {
		view: LivePadView;
		canFinish: boolean;
		onScore: (code: ScoreCode) => void;
		onFinish: () => void;
	} = $props();
</script>

<p>{view.progress}</p>
<p class="total">{view.total}</p>
<p>{summaryLine({ total: view.total, hits: view.hits, golds: view.golds, xCount: view.xCount, arrowCount: 0 })}</p>
<ol class="slots">
	{#each view.slots as slot, i (i)}
		<li>{slot ?? ''}</li>
	{/each}
</ol>
<div class="pad">
	{#each pad as row, r (r)}
		<div class="row">
			{#each row as code (code)}
				<button class="hit" type="button" onclick={() => onScore(code)}>{code}</button>
			{/each}
		</div>
	{/each}
</div>
<button class="hit primary" type="button" disabled={!canFinish} onclick={onFinish}>Finish session</button>
