<script lang="ts">
	import type { Archer } from '$lib/api/types';

	let { archers, onPick }: { archers: Archer[]; onPick: (id: string) => void } = $props();
</script>

<h1>Home Archery</h1>
<p>Who is shooting?</p>
{#each archers as archer (archer.id)}
	<button class="hit" type="button" onclick={() => onPick(archer.id)}>{archer.displayName}</button>
{/each}
