export const CURRENT_SESSION_KEY = 'currentSessionId';

export function getCurrentSessionId(): string | null {
	if (typeof sessionStorage === 'undefined') {
		return null;
	}
	return sessionStorage.getItem(CURRENT_SESSION_KEY);
}

export function setCurrentSessionId(id: string): void {
	sessionStorage.setItem(CURRENT_SESSION_KEY, id);
}

export function clearCurrentSessionId(): void {
	sessionStorage.removeItem(CURRENT_SESSION_KEY);
}
