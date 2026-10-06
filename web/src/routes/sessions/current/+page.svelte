<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import LivePad from '$lib/components/LivePad.svelte';
	import { api } from '$lib/api/client';
	import { getCurrentSessionId, clearCurrentSessionId } from '$lib/session/current';
	import { livePadView, type LivePadView } from '$lib/session/mappers';
	import type { ScoreCode, Session } from '$lib/api/types';

	let view = $state<LivePadView | null>(null);
	let canFinish = $state(false);
	let sessionId = $state<string | null>(null);

	async function load(id: string) {
		const session: Session = await api.getSession(id);
		view = livePadView(session);
		canFinish = session.summary.arrowCount > 0 && session.status === 'in_progress';
	}

	onMount(async () => {
		const id = getCurrentSessionId();
		if (!id) {
			await goto('/home');
			return;
		}
		sessionId = id;
		await load(id);
	});

	async function score(code: ScoreCode) {
		if (!sessionId) {
			return;
		}
		await api.recordArrow(sessionId, code);
		await load(sessionId);
	}

	async function finish() {
		if (!sessionId) {
			return;
		}
		await api.completeSession(sessionId);
		clearCurrentSessionId();
		await goto('/home');
	}
</script>

<svelte:head>
	<title>Live session</title>
</svelte:head>

{#if view}
	<LivePad {view} {canFinish} onScore={score} onFinish={finish} />
{/if}
