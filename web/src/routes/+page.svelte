<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import ArcherPicker from '$lib/components/ArcherPicker.svelte';
	import { api } from '$lib/api/client';
	import { setArcherId } from '$lib/session/archer';
	import type { Archer } from '$lib/api/types';

	let archers = $state<Archer[]>([]);

	onMount(async () => {
		archers = await api.listArchers();
	});

	function pick(id: string) {
		setArcherId(id);
		void goto('/home');
	}
</script>

<svelte:head>
	<title>Who is shooting?</title>
</svelte:head>

<ArcherPicker {archers} onPick={pick} />
