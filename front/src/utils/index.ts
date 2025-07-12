export function debounce(fn: (...args: any[]) => void, delay: number) {
    let timer: number | null = null;
    return (...args: any[]) => {
        if (timer) {
            clearTimeout(timer);
        }
        timer = setTimeout(() => fn(...args), delay);
    };
}