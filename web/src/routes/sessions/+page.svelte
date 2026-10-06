<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import HistoryList from '$lib/components/HistoryList.svelte';
	import { api } from '$lib/api/client';
	import { getArcherId } from '$lib/session/archer';
	import { sessionListRows, type SessionListRow } from '$lib/session/mappers';

	let sessions = $state<SessionListRow[]>([]);

	onMount(async () => {
		const archerId = getArcherId();
		if (!archerId) {
			await goto('/');
			return;
		}
		const [list, bests] = await Promise.all([
			api.listSessions(archerId),
			api.getPersonalBests(archerId)
		]);
		sessions = sessionListRows(list, bests);
	});
</script>

<svelte:head>
	<title>Sessions</title>
</svelte:head>

<HistoryList
	{sessions}
	onOpen={(id) => goto(`/sessions/${id}`)}
	onHome={() => goto('/home')}
/>
