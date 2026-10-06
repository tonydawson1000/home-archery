import { beforeEach, describe, expect, it } from 'vitest';
import { ARCHER_ID_KEY, clearArcherId, getArcherId, setArcherId } from '$lib/session/archer';
import { TONY_ID } from '$lib/api/types';

const memory = new Map<string, string>();

beforeEach(() => {
	memory.clear();
	const localStorage = {
		getItem: (key: string) => memory.get(key) ?? null,
		setItem: (key: string, value: string) => {
			memory.set(key, value);
		},
		removeItem: (key: string) => {
			memory.delete(key);
		},
		clear: () => memory.clear()
	};
	Object.defineProperty(globalThis, 'localStorage', { value: localStorage, configurable: true });
});

describe('archer localStorage', () => {
	it('round-trips archerId', () => {
		expect(getArcherId()).toBeNull();
		setArcherId(TONY_ID);
		expect(localStorage.getItem(ARCHER_ID_KEY)).toBe(TONY_ID);
		expect(getArcherId()).toBe(TONY_ID);
		clearArcherId();
		expect(getArcherId()).toBeNull();
	});
});
