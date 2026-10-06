<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import HomeScreen from '$lib/components/HomeScreen.svelte';
	import { api } from '$lib/api/client';
	import { clearArcherId, getArcherId } from '$lib/session/archer';
	import { setCurrentSessionId } from '$lib/session/current';
	import { lastFive, personalBestCard, sessionListRows, type SessionListRow } from '$lib/session/mappers';
	import type { Archer } from '$lib/api/types';
	import type { PersonalBestCard } from '$lib/session/mappers';

	let displayName = $state('');
	let pb = $state<PersonalBestCard>({ headline: 'No personal best yet', hint: null });
	let sessions = $state<SessionListRow[]>([]);

	onMount(async () => {
		const archerId = getArcherId();
		if (!archerId) {
			await goto('/');
			return;
		}
		const [archers, list, bests] = await Promise.all([
			api.listArchers(),
			api.listSessions(archerId),
			api.getPersonalBests(archerId)
		]);
		displayName = archers.find((a: Archer) => a.id === archerId)?.displayName ?? '';
		pb = personalBestCard(bests);
		sessions = lastFive(sessionListRows(list, bests));
	});

	async function start() {
		const archerId = getArcherId();
		if (!archerId) {
			return;
		}
		const session = await api.startSession(archerId);
		setCurrentSessionId(session.id);
		await goto('/sessions/current');
	}
</script>

<svelte:head>
	<title>Home</title>
</svelte:head>

<HomeScreen
	{displayName}
	{pb}
	{sessions}
	onStart={start}
	onSwitch={() => {
		clearArcherId();
		void goto('/');
	}}
	onOpenSession={(id) => goto(`/sessions/${id}`)}
	onHistory={() => goto('/sessions')}
/>
