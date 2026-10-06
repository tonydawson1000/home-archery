export const ARCHER_ID_KEY = 'archerId';

export function getArcherId(): string | null {
	if (typeof localStorage === 'undefined') {
		return null;
	}
	return localStorage.getItem(ARCHER_ID_KEY);
}

export function setArcherId(id: string): void {
	localStorage.setItem(ARCHER_ID_KEY, id);
}

export function clearArcherId(): void {
	localStorage.removeItem(ARCHER_ID_KEY);
}
