import { SearchResult } from '../models/search-result.model';

export class ApiError extends Error {
    constructor(message: string) {
        super(message);
    }
}

export class SearchService {
    private static BASE_URL = 'http://127.0.0.1:8080';

    static async search(value: string): Promise<SearchResult> {
        if (!value) {
            throw new ApiError('Invalid input: Value cannot be empty');
        }

        const sanitizedValue = this.sanitizeInput(value);

        try {
            const response = await fetch(`${this.BASE_URL}/numbers/${encodeURIComponent(sanitizedValue)}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({}));
                throw new ApiError(errorData.message || `HTTP error! status: ${response.status}`);
            }

            const data: SearchResult = await response.json();

            if (!this.isValidSearchResult(data)) {
                throw new ApiError('Invalid response format');
            }

            return data;
        } catch (error) {
            console.error('Search API Error:', error);

            if (error instanceof ApiError) {
                throw error;
            }

            throw new ApiError('An unexpected error occurred during the search');
        }
    }

    private static sanitizeInput(input: string): string {
        return input
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .trim();
    }

    // generally prefer to avoid any but in this case seems valid for a guard
    private static isValidSearchResult(result: any): result is SearchResult {
        return (
            typeof result === 'object' &&
            typeof result.index === 'number' &&
            typeof result.value === 'number' &&
            (result.message === undefined || typeof result.message === 'string')
        );
    }
}