
import { SearchService } from '../services/search.service';
import { ApiError } from '../services/search.service';

export class SearchUI {
    private valueInput: HTMLInputElement;
    private resultMessage: HTMLParagraphElement;
    private searchBtn: HTMLButtonElement;

    constructor() {
        this.valueInput = this.getElement<HTMLInputElement>('valueInput');
        this.resultMessage = this.getElement<HTMLParagraphElement>('resultMessage');
        this.searchBtn = this.getElement<HTMLButtonElement>('searchBtn');

        this.bindEvents();
    }

    private getElement<T extends HTMLElement>(id: string): T {
        const element = document.getElementById(id);
        if (!element) {
            throw new Error(`Element with id ${id} not found`);
        }
        return element as T;
    }

    private bindEvents() {
        this.searchBtn.addEventListener('click', () => this.submitSearch());
    }

    private async submitSearch() {
        const value = this.valueInput.value;

        this.resultMessage.textContent = '';
        this.resultMessage.className = '';

        try {
            const result = await SearchService.search(value);

            this.resultMessage.className = 'success';
            this.resultMessage.textContent =
                `Index: ${result.index}, Value: ${result.value}, ` +
                `Message: ${result.message || 'No additional message'}`;
        } catch (error) {
            this.resultMessage.className = 'error';
            this.resultMessage.textContent =
                error instanceof ApiError
                    ? error.message
                    : 'An unexpected error occurred';
        }
    }
}