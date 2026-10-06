<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import Scorecard from '$lib/components/Scorecard.svelte';
	import { api } from '$lib/api/client';
	import { setCurrentSessionId } from '$lib/session/current';
	import type { Session } from '$lib/api/types';

	let session = $state<Session | null>(null);

	onMount(async () => {
		const id = page.params.id;
		if (!id) {
			return;
		}
		session = await api.getSession(id);
	});

	function resume() {
		if (!session) {
			return;
		}
		setCurrentSessionId(session.id);
		void goto('/sessions/current');
	}
</script>

<svelte:head>
	<title>Session</title>
</svelte:head>

{#if session}
	<Scorecard {session} onResume={resume} />
{/if}
